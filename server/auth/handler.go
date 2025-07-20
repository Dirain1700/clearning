// Package auth provides authentication handlers for the application.
package auth

import (
	"fmt"
	"net/http"

	"github.com/dirain1700/clearning/server/def"
	"github.com/markbates/goth/gothic"
)

// HandleAuthEntry processes the authentication entry point.
func HandleAuthEntry(res http.ResponseWriter, req *http.Request) {
	// try to get the user without re-authenticating
	def.ModifyURLParams(req, "provider", "google")
	_, err := gothic.CompleteUserAuth(res, req)
	fmt.Printf("Attempting to complete user auth: %v\n", err)

	if err == nil {
		http.Redirect(res, req, "/info", http.StatusFound)
	} else {
		gothic.BeginAuthHandler(res, req)
	}
}
