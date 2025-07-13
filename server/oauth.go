// Package server (Oauth) provides the OAuth functionality for the application.
package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/joho/godotenv"
)

// User represents a user in the system.
type User struct {
	id string
}

func getGoogleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       []string{"profile", "openid"},
		Endpoint:     google.Endpoint,
	}
}

func getGoogleOauthURL(config *oauth2.Config) string {
	return config.AuthCodeURL("state")
}

func getUserInfo(ctx context.Context, config *oauth2.Config, token *oauth2.Token) (*User, error) {
	client := config.Client(ctx, token)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	defer func() {
		err = res.Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body: %v", err)
		}
	}()

	userInfo := make(map[string]interface{})

	err = json.NewDecoder(res.Body).Decode(&userInfo)
	if err != nil {
		return &User{id: ""}, fmt.Errorf("failed to decode user info: %w", err)
	}

	id, ok := userInfo["id"].(string)
	if !ok {
		return nil, errors.New("userInfo[\"id\"] is not a string") //nolint:err113 // This is a custom error message for clarity.
	}

	return &User{id: id}, nil
}

func oauth(writer http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")

	config := getGoogleOauthConfig()

	token, err := config.Exchange(req.Context(), code)
	if err != nil {
		http.Error(writer, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)

		return
	}

	user, err := getUserInfo(req.Context(), config, token)
	if err != nil {
		fmt.Println("Failed to get user info:", err)

		return
	}

	fmt.Println("Log in successful")

	_, err = writer.Write([]byte(user.id))
	if err != nil {
		http.Error(writer, "Failed to write response: "+err.Error(), http.StatusInternalServerError)

		return
	}
}

// TestOauth is just a test function to run the OAuth flow.
func TestOauth() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")

		return
	}

	fmt.Println(getGoogleOauthURL(getGoogleOauthConfig()))
	http.HandleFunc("/callback", oauth)

	err = http.ListenAndServe(":8000", nil) //nolint:gosec // This is just a test server, not for production use.
	if err != nil {
		fmt.Println("Error starting server:", err)

		return
	}
}
