package jwt

import (
	"errors"
	"log"
	"net/http"
)

func SetCookieHandler(w http.ResponseWriter, r *http.Request, token string, duration int, name string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		MaxAge:   duration,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
}

func GetCookieHandler(w http.ResponseWriter, r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return ""
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return ""
	}
	return cookie.Value
}

func DeleteCookieHandler(w http.ResponseWriter, r *http.Request, name string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}
