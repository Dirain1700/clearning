// Package logout provides the logout handler for the application.
package logout

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dirain1700/clearning/server/def"
	"github.com/markbates/goth/gothic"
)

// HandleLogout processes the logout request and clears the JWT cookie.
func HandleLogout(res http.ResponseWriter, req *http.Request) {
	def.ModifyURLParams(req, "provider", "google")

	err := gothic.Logout(res, req)
	if err != nil {
		log.Fatalln(fmt.Errorf("error logging out: %w", err))
	}

	res.Header().Set("Location", "/")
	http.SetCookie(res, &http.Cookie{
		Name:     def.JWTCookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
		MaxAge:   -1, // Expire the cookie immediately
	})

	res.WriteHeader(http.StatusOK)
}
