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

func IsTokenInvalid(tokenString string, secretKey string) (bool, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
	})
	if err != nil {
			return false, err
	}

	// Check if token is valid and contains expiration claim
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Retrieve expiration time
			exp := int64(claims["exp"].(float64))

			// Compare expiration time with current time
			if exp < time.Now().Unix() {
					// Token has expired
					return true, nil
			}

			// Token is not expired
			return false, nil
	}

	// Token is invalid or does not contain expiration claim
	return false, nil
}
