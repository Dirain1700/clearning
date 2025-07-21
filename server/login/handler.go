// Package login provides the login handler for the application.
package login

import (
	"log"
	"net/http"
	"net/url"

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
		log.Printf("error completing user auth: %v", err)
	}

	jwt, err := auth.GenerateJWT(
		&goth.User{
			UserID:    user.UserID,
			Name:      user.Name,
			AvatarURL: user.AvatarURL,
		},
	)
	if err != nil {
		log.Printf("error generating JWT: %v", err)

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

	encodedPathFrom := req.URL.Query().Get("from")

	pathFrom, err := url.PathUnescape(encodedPathFrom)
	if err != nil {
		log.Printf("Error decoding 'from' parameter: %v", err)

		pathFrom = "/" // Default redirect path if decoding fails
	}

	if pathFrom == "" {
		pathFrom = "/" // Default redirect path if not specified
	}

	http.Redirect(res, req, pathFrom, http.StatusFound)
}
