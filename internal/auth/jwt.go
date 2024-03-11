package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateJWT(user User, jwtSecret string, expiresIn ...int) string {
	maxAge := time.Now().Add(24 * time.Hour)
	expIn := maxAge

	// Set expire time to optional expiresIn time, as long as it is less than 24 hours
	if len(expiresIn) > 0 && maxAge.Sub(expIn) > 0 {
		expIn = time.Now().Add(time.Duration(expiresIn[0]) * time.Second)
	}
	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "chirpy",
		ExpiresAt: jwt.NewNumericDate(expIn),
		Subject:   fmt.Sprint(user.Id),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := jwtToken.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Fatal(err)
	}

	return ss
}

// TODO: function
func ValidateJWT(header string, email, password string) (bool, error) {

	return true, nil
}
