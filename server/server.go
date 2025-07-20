package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"

	"github.com/dirain1700/clearning/server/auth"
	"github.com/dirain1700/clearning/server/info"
	"github.com/dirain1700/clearning/server/login"
	"github.com/dirain1700/clearning/server/logout"
)

// TestProvider initializes the provider and sets up routes for authentication.
func TestProvider() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")

		return
	}

	key := []byte(os.Getenv("SESSION_SECRET"))
	cookieStore := sessions.NewCookieStore(key)
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(auth.JWTExpiresIn.Seconds()),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteDefaultMode,
	}
	gothic.Store = cookieStore

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), os.Getenv("REDIRECT_URL"), "profile"),
	)

	http.HandleFunc("/api/auth/callback", login.HandleLogin)

	http.HandleFunc("/api/logout", logout.HandleLogout)

	http.HandleFunc("/api/auth", auth.HandleAuthEntry)

	http.HandleFunc("/api/info", info.HandleInformation)

	log.Println("listening on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil)) //nolint: gosec // This is just a test server, not for production use.
}
