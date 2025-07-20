// Package def provides utility constants and functions for the application.
package def

import "net/http"

const (
	// JWTCookieName defines the name of the cookie used for authentication.
	JWTCookieName = "google_auth"
)

// ModifyURLParams modifies the URL query parameters of the given request.
func ModifyURLParams(req *http.Request, key, value string) {
	// Modify the URL query parameters
	query := req.URL.Query()
	query.Set(key, value)
	req.URL.RawQuery = query.Encode()
}
