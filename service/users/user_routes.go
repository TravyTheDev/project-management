package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-management/service/i18n"
	"project-management/service/jwt"
	"project-management/types"
	"project-management/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userStore    types.UserStore
	sessionStore types.SessionStore
	jwtMaker     *jwt.JWTMaker
	langMap      i18n.LangMap
}

func NewHandler(userStore types.UserStore, sessionStore types.SessionStore, secretKey string, langMap i18n.LangMap) *UserHandler {
	return &UserHandler{
		userStore:    userStore,
		sessionStore: sessionStore,
		jwtMaker:     jwt.NewJWTMaker(secretKey),
		langMap:      langMap,
	}
}

func (h *UserHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register/{lang}", h.handleRegister).Methods("POST")
	router.HandleFunc("/login/{lang}", h.handleLogin).Methods("POST")
	router.HandleFunc("/logout", h.handleLogout).Methods("POST")
	router.HandleFunc("/me", jwt.GetAuthMiddlewareFunc(h.jwtMaker, h.handleGetUser)).Methods("GET")
	router.HandleFunc("/renew_token", h.renewAccessToken).Methods("POST")
	router.HandleFunc("/search_user/{email}", h.handleSearchByEmail).Methods("GET")
	router.HandleFunc("/confirm_numbers", h.confirmNumbers).Methods("POST")
	router.HandleFunc("/change_password", h.changePassword).Methods("POST")
}

func (h *UserHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	lang := mux.Vars(r)["lang"]
	var payload types.RegisterPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Println(err)
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		if errors != nil {
			http.Error(w, h.langMap[lang]["auth"]["invalid_payload"], http.StatusBadRequest)
		}
		return
	}

	if payload.Password != payload.PasswordConfirm {
		http.Error(w, h.langMap[lang]["auth"]["password_mismatch"], http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		fmt.Println(err)
	}

	user := types.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashedPassword,
		IsAdmin:  false,
	}

	u, err := h.userStore.CreateUser(user)
	if err != nil {
		http.Error(w, h.langMap[lang]["auth"]["email_in_use"], http.StatusConflict)
		return
	}
	if u.ID == 0 {
		fmt.Println("error creating user")
		return
	}

	h.setCookiesAndSession(u, w, r)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	lang := mux.Vars(r)["lang"]
	var payload types.LoginUserPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Println(err)
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		if errors != nil {
			http.Error(w, h.langMap[lang]["auth"]["no_email_no_pass"], http.StatusBadRequest)
		}
		return
	}

	u, err := h.userStore.GetUserByEmail(payload.Email)
	if err != nil {
		http.Error(w, h.langMap[lang]["auth"]["invalid"], http.StatusBadRequest)
		return
	}
	if err := utils.CheckPassword(payload.Password, u.Password); err != nil {
		http.Error(w, h.langMap[lang]["auth"]["invalid"], http.StatusBadRequest)
		return
	}

	h.setCookiesAndSession(u, w, r)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func (h *UserHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie := jwt.GetCookieHandler(w, r, "refresh")
	claims, err := h.jwtMaker.VerifyToken(cookie)
	if err != nil {
		http.Error(w, "error verifying token", http.StatusInternalServerError)
		return
	}

	if err := h.sessionStore.DeleteSession(claims.RegisteredClaims.ID); err != nil {
		http.Error(w, "error deleting session", http.StatusInternalServerError)
		return
	}
	jwt.DeleteCookieHandler(w, r, "authentication")
	jwt.DeleteCookieHandler(w, r, "refresh")
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value(jwt.AuthKey{}).(*jwt.UserClaims)

	userID := userClaims.ID

	user, err := h.userStore.GetUserByID(userID)
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) handleSearchByEmail(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]
	u, err := h.userStore.GetUserByEmail(email)
	if err != nil {
		return
	}

	user := &types.UserRes{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return
	}
}

func (h *UserHandler) renewAccessToken(w http.ResponseWriter, r *http.Request) {

	refreshToken := jwt.GetCookieHandler(w, r, "refresh")

	refreshClaims, err := h.jwtMaker.VerifyToken(refreshToken)
	if err != nil {
		http.Error(w, "error verifying token", http.StatusUnauthorized)
		return
	}

	session, err := h.sessionStore.GetSession(refreshClaims.RegisteredClaims.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error getting sessions", http.StatusInternalServerError)
		return
	}

	if session.IsRevoked {
		http.Error(w, "session revoked", http.StatusUnauthorized)
		return
	}

	if session.UserEmail != refreshClaims.Email {
		http.Error(w, "session revoked", http.StatusUnauthorized)
		return
	}

	var duration = time.Minute * 15

	accessToken, _, err := h.jwtMaker.CreateToken(
		refreshClaims.ID,
		refreshClaims.Username,
		refreshClaims.Email,
		refreshClaims.IsAdmin,
		duration,
	)
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}

	jwt.SetCookieHandler(w, r, accessToken, int(duration.Seconds()), "authentication")

}

func (h *UserHandler) setCookiesAndSession(u *types.User, w http.ResponseWriter, r *http.Request) {
	accessDuration := time.Minute * 15
	refreshDuration := 24 * time.Hour

	accessToken, _, err := h.jwtMaker.CreateToken(int(u.ID), u.Username, u.Email, u.IsAdmin, accessDuration)
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}

	refreshToken, refreshClaims, err := h.jwtMaker.CreateToken(int(u.ID), u.Username, u.Email, u.IsAdmin, refreshDuration)
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}

	sesh := &types.Session{
		ID:           refreshClaims.RegisteredClaims.ID,
		UserEmail:    u.Email,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		CreatedAt:    time.Now(),
		ExpiresAt:    refreshClaims.ExpiresAt.Time,
	}

	_, err = h.sessionStore.CreateSession(sesh)
	if err != nil {
		http.Error(w, "error creating session", http.StatusInternalServerError)
		return
	}

	jwt.SetCookieHandler(w, r, accessToken, int(accessDuration.Seconds()), "authentication")
	jwt.SetCookieHandler(w, r, refreshToken, int(refreshDuration.Seconds()), "refresh")
}

func (h *UserHandler) confirmNumbers(w http.ResponseWriter, r *http.Request) {
	var numbers types.NumbersConfirmReq
	cookie := jwt.GetCookieHandler(w, r, "password_reset")

	if err := json.NewDecoder(r.Body).Decode(&numbers); err != nil {
		http.Error(w, "error decoding numbers", http.StatusInternalServerError)
		return
	}

	claims, err := h.jwtMaker.VerifyPasswordResetToken(cookie)
	if err != nil {
		http.Error(w, "error verifying token", http.StatusInternalServerError)
		return
	}

	if claims.Numbers != numbers.Numbers {
		http.Error(w, "numbers don't match", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) changePassword(w http.ResponseWriter, r *http.Request) {
	var req types.PasswordResetReq
	cookie := jwt.GetCookieHandler(w, r, "password_reset")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding numbers", http.StatusInternalServerError)
		return
	}

	claims, err := h.jwtMaker.VerifyPasswordResetToken(cookie)
	if err != nil {
		http.Error(w, "error verifying token", http.StatusInternalServerError)
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		errors := err.(validator.ValidationErrors)
		if errors != nil {
			http.Error(w, "email must be valid\npassword must be\nat least 3 characters", http.StatusBadRequest)
		}
		return
	}

	if req.Password != req.PasswordConfirm {
		http.Error(w, "passwords must match", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		fmt.Println(err)
	}

	if err := h.userStore.ChangePassword(claims.Email, hashedPassword); err != nil {
		http.Error(w, "error changing password", http.StatusBadRequest)
		return
	}
}
