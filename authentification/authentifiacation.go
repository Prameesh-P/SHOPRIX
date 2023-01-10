package authentification

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

var JwtKey = []byte(os.Getenv("SUPER_SECRET"))

type JWTCliam struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(email string) (map[string]string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTCliam{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenstring, err := token.SignedString(JwtKey)
	if err != nil {
		return nil, err
	}
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtCliams := refreshToken.Claims.(jwt.MapClaims)
	rtCliams["email"] = email
	rtCliams["exp"] = time.Now().Add(24 * time.Hour)

	rt, err := refreshToken.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  tokenstring,
		"refresh_token": rt,
	}, nil
}

var P string

func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTCliam{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(JwtKey), nil
		})
	if err != nil {
		return
	}
	cliams, Ok := token.Claims.(*JWTCliam)
	P = cliams.Email
	if !Ok {
		err = errors.New("could not parse claims")
		return
	}
	if cliams.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return
}
