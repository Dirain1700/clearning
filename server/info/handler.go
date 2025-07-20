// Package info provides the handler for retrieving server information.
package info

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dirain1700/clearning/server/auth"
	"github.com/dirain1700/clearning/server/def"
)

// UserInfo represents the structure of user information returned by the server.
type UserInfo struct {
	Subject   string `json:"sub"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatarUrl"`
}

// HandleInformation processes the request to retrieve server information.
func HandleInformation(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	jwtCookie, err := req.Cookie(def.JWTCookieName)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)

		err = json.NewEncoder(res).Encode(&def.ErrorResponse{
			Message: "Could not read JWT from cookie",
		})
		if err != nil {
			log.Printf("Failed to encode error response: %v", err)
		}

		return
	}

	user, err := auth.VerifyJWT(jwtCookie.Value)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)

		err = json.NewEncoder(res).Encode(&def.ErrorResponse{
			Message: "Invalid JWT token",
		})
		if err != nil {
			log.Printf("Failed to encode error response: %v", err)
		}

		return
	}

	jsonResponse := &UserInfo{
		Subject:   user.Subject,
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
	}

	err = json.NewEncoder(res).Encode(jsonResponse)
	if err != nil {
		log.Printf("Failed to encode user information response: %v", err)
	}
}
