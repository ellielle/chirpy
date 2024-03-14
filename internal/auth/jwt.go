package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type Claims struct {
	jwt.RegisteredClaims
}

func CreateJWT(user User, jwtSecret string, isAccess bool) (string, error) {
	// Set max age for access tokens to 1 hour
	maxAge := time.Now().Add(1 * time.Hour)
	issuer := "chirpy-access"

	// Set a 60day expiration and change issuer if the token is meant
	// to be a refresh token
	if !isAccess {
		maxAge = time.Now().Add(24 * 60 * time.Hour)
		issuer = "chirpy-refresh"
	}

	claims := &Claims{
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(maxAge),
			Subject:   fmt.Sprint(user.Id),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := jwtToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return ss, nil
}

func ValidateJWT(token, jwtSecret string) (*jwt.Token, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return &jwt.Token{}, err
	}

	return jwtToken, nil
}

func GetUserIDWithToken(token jwt.Token) (string, error) {
	userID, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}
	return userID, nil
}

func GetJWTIssuer(token, jwtSecret string) (string, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", err
	}

	issuer, err := jwtToken.Claims.GetIssuer()
	if err != nil {
		return "", err
	}

	return issuer, nil
}

func RefreshJWT(token, jwtSecret string) (string, error) {

	return "not finished", nil
}
