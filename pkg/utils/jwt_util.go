package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateAccessToken(id string) (string, error) {
	config, err := ReadAppConfig()
	if err != nil {
		return "cant read config", err
	}


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	})
	return token.SignedString([]byte(config.Jwt.AtSecret))
}

func GenerateRefreshToken(id string) (string, error) {
	config, err := ReadAppConfig()
	if err != nil {
		return "cant read config", err
	}


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	return token.SignedString([]byte(config.Jwt.RtSecret))
}
