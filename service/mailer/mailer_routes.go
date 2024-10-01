package mailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"project-management/service/jwt"
	"project-management/types"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/wneessen/go-mail"
)

type MailHandler struct {
	userStore types.UserStore
	jwtMaker  *jwt.JWTMaker
}

func NewMailHandler(userStore types.UserStore, secretKey string) *MailHandler {
	return &MailHandler{
		userStore: userStore,
		jwtMaker:  jwt.NewJWTMaker(secretKey),
	}
}

func (m *MailHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/mailer/forgot_password", m.sendForgotPasswordEmail).Methods("POST")
}

func (m *MailHandler) sendForgotPasswordEmail(w http.ResponseWriter, r *http.Request) {
	var req types.ForgotPasswordReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error getting email", http.StatusBadRequest)
		return
	}
	user, err := m.userStore.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	username := os.Getenv("MAILTRAP_USER")
	password := os.Getenv("MAILTRAP_PASS")
	smtpHost := os.Getenv("MAILTRAP_HOST")
	smtpPort := os.Getenv("MAILTRAP_PORT")
	from := os.Getenv("SENDER")

	port, _ := strconv.Atoi(smtpPort)

	t, _ := template.ParseFiles("service/mailer/forgot_password_template.html")

	var body bytes.Buffer
	var numbers = rangeIn(10000, 99999)
	if err := t.Execute(&body, struct {
		Email   string
		Numbers int
	}{
		Email:   user.Email,
		Numbers: numbers,
	}); err != nil {
		fmt.Println(err)
	}

	mailer := mail.NewMsg()
	if err := mailer.From(from); err != nil {
		log.Fatalf("failed to set From address: %s", err)
	}
	if err := mailer.To(req.Email); err != nil {
		log.Fatalf("failed to set To address: %s", err)
	}
	mailer.Subject("Password reset")
	mailer.SetBodyString(mail.TypeTextHTML, body.String())

	c, err := mail.NewClient(smtpHost,
		mail.WithPort(port), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(username), mail.WithPassword(password))
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}

	if err := c.DialAndSend(mailer); err != nil {
		log.Fatalf("failed to send mail: %s", err)
	}

	passwordResetToken, _, err := m.jwtMaker.CreatePasswordResetToken(req.Email, numbers, time.Minute*5)
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}
	jwt.SetCookieHandler(w, r, passwordResetToken, int(time.Minute*5), "password_reset")
}

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}
