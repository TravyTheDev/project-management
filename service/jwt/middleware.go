package jwt

import (
	"context"
	"fmt"
	"net/http"
)

type AuthKey struct{}

func GetAuthMiddlewareFunc(jwtMaker *JWTMaker, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := verifyClaimsFromCookie(w, r, jwtMaker)
		if err != nil {
			http.Error(w, fmt.Sprintf("error verifying token: %v", err), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), AuthKey{}, claims)
		handlerFunc(w, r.WithContext(ctx))
	}
}

func verifyClaimsFromCookie(w http.ResponseWriter, r *http.Request, jwtMaker *JWTMaker) (*UserClaims, error) {
	cookie := GetCookieHandler(w, r, "authentication")
	claims, err := jwtMaker.VerifyToken(cookie)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return claims, nil
}
