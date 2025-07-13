package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/joho/godotenv"
)

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

func getGoogleOauthUrl(config *oauth2.Config) string {
	return config.AuthCodeURL("state")
}

func getUserInfo(config *oauth2.Config, token *oauth2.Token) (*User, error) {
	client := config.Client(context.TODO(), token)

	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}

	userInfo := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&userInfo)
	if err != nil {
		return &User{id: ""}, err
	}

	id := userInfo["id"].(string)
	return &User{id: id}, nil
}

func oauth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	config := getGoogleOauthConfig()
	token, err := config.Exchange(context.TODO(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := getUserInfo(config, token)
	if err != nil {
		fmt.Println("Failed to get user info:", err)
		return
	}

	fmt.Println("Log in successful")
	w.Write([]byte(user.id))
}

func TestOauth() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	fmt.Println(getGoogleOauthUrl(getGoogleOauthConfig()))
	http.HandleFunc("/callback", oauth)
	http.ListenAndServe(":8000", nil)
}
