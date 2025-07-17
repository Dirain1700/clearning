package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func afterLoginHandler(res http.ResponseWriter, req *http.Request) {
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		log.Fatalln(fmt.Errorf("error completing user auth: %w", err))
	}

	jwt, err := GenerateJWT(
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
		Name:     "google_auth",
		Value:    jwt,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
		MaxAge:   int(jwtExpiresIn),
	})

	_, err = fmt.Fprintln(res, "Login successful! JWT set in cookie. Your id:", user.UserID)
	if err != nil {
		log.Fatalln(fmt.Errorf("error writing response: %w", err))
	}
}

func logoutHandler(res http.ResponseWriter, req *http.Request) {
	err := gothic.Logout(res, req)
	if err != nil {
		log.Fatalln(fmt.Errorf("error logging out: %w", err))
	}

	res.Header().Set("Location", "/")
	res.WriteHeader(http.StatusTemporaryRedirect)
	http.SetCookie(res, &http.Cookie{
		Name:     "google_auth",
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
		MaxAge:   int(jwtExpiresIn),
	})
}

func authHandler(res http.ResponseWriter, req *http.Request) {
	// try to get the user without re-authenticating
	gothUser, err := gothic.CompleteUserAuth(res, req)
	if err == nil {
		_, err = fmt.Fprintln(res, "Found user:", gothUser.UserID)
		if err != nil {
			log.Fatalln(fmt.Errorf("error writing response: %w", err))
		}
	} else {
		gothic.BeginAuthHandler(res, req)
	}
}

func infoHandler(res http.ResponseWriter, req *http.Request) {
	jwtCookie, err := req.Cookie("google_auth")
	if err == nil {
		user, err := VerifyJWT(jwtCookie.Value)
		if err != nil {
			// fmt.Fprintln(res, "Error verifying JWT:", err)
			http.Redirect(res, req, "/auth/google", http.StatusFound)

			return
		}

		_, err = fmt.Fprintf(res, "Welcome back!\nSubject: %q\nName: %q\nAvatarURL: %q\nExp: %d\n",
			user.Subject, user.Name, user.AvatarURL, user.Exp)
		if err != nil {
			log.Fatalln(fmt.Errorf("error writing response: %w", err))
		}
	} else {
		_, err = fmt.Fprintln(res, "Please log in using one of the providers.")
		if err != nil {
			log.Fatalln(fmt.Errorf("error writing response: %w", err))
		}
	}
}

// TestProvider initializes the provider and sets up routes for authentication.
func TestProvider() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")

		return
	}

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), os.Getenv("REDIRECT_URL"), "profile"),
	)

	mapPro := map[string]string{
		"google": "Google",
	}

	keys := make([]string, 0, len(mapPro))
	for k := range mapPro {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	http.HandleFunc("/auth/{provider}/callback", afterLoginHandler)

	http.HandleFunc("/logout/{provider}", logoutHandler)

	http.HandleFunc("/auth/{provider}", authHandler)

	http.HandleFunc("/info", infoHandler)

	http.HandleFunc("/", infoHandler)

	log.Println("listening on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil)) //nolint: gosec // This is just a test server, not for production use.
}

// ProviderIndex holds the list of providers and their names.
type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}
