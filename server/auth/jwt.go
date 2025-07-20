// Package auth provides JWT auth func for the application.
package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/markbates/goth"
)

const (
	rsaKeyPath       = "./server/rsa.key"
	rsaPublicKeyPath = "./server/rsa_pub.pem"

	daysInWeek = 7
	hoursInDay = 24

	// JWTExpiresIn defines the expiration time for the JWT token.
	JWTExpiresIn = (time.Hour * hoursInDay * daysInWeek) // 7 days
)

var errSignMethosdMismatch = errors.New("signing method mismatch")

// UserJWTClaims defines the structure of the JWT claims for the user.
type UserJWTClaims struct {
	jwt.StandardClaims

	Name      string `json:"name"`
	AvatarURL string `json:"avatarUrl"`
	Exp       int64  `json:"exp"`
}

// GenerateJWT generates a JWT for the user.
func GenerateJWT(user *goth.User) (string, error) {
	privateKeyStr, err := os.ReadFile(rsaKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create the JWT claims and sign
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &UserJWTClaims{
		StandardClaims: jwt.StandardClaims{
			Subject: user.UserID,
			// No need to use Seconds() cuz JWTExpiresIn is already in time.Duration format
			ExpiresAt: jwt.TimeFunc().Add(JWTExpiresIn).Unix(),
		},
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
		// No need to use Seconds() cuz JWRExpiresIn is already in time.Duration format
		Exp: jwt.TimeFunc().Add(JWTExpiresIn).Unix(),
	})

	tokenStr, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenStr, nil
}

// VerifyJWT verifies the JWT token and returns the user information.
func VerifyJWT(tokenStr string) (*UserJWTClaims, error) {
	publicKeyStr, err := os.ReadFile(rsaPublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	token, err := jwt.ParseWithClaims(tokenStr, &UserJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("%w: %v", errSignMethosdMismatch, token.Header["alg"])
		}

		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*UserJWTClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
