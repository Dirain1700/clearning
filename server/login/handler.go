// Package login provides the login handler for the application.
package login

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dirain1700/clearning/server/auth"
	"github.com/dirain1700/clearning/server/def"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

// HandleLogin processes the login request and generates a JWT for the user.
func HandleLogin(res http.ResponseWriter, req *http.Request) {
	def.ModifyURLParams(req, "provider", "google")

	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		log.Fatalln(fmt.Errorf("error completing user auth: %w", err))
	}

	jwt, err := auth.GenerateJWT(
		&goth.User{
			UserID:    user.UserID,
			Name:      user.Name,
			AvatarURL: user.AvatarURL,
		},
	)
	if err != nil {
		log.Fatalln(fmt.Errorf("error generating JWT: %w", err))

		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:     def.JWTCookieName,
		Value:    jwt,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
		MaxAge:   int(auth.JWTExpiresIn.Seconds()),
	})

	res.WriteHeader(http.StatusOK)
}
