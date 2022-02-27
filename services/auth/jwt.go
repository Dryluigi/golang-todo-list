package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTContextKey string

type JwtClaims struct {
	UserId uint `json:"user_id"`
	jwt.StandardClaims
}

type JwtClaimsContent struct {
	UserId uint   `json:"user_id"`
	Exp    uint64 `json:"exp"`
}

type IssuedToken struct {
	AccessToken string `json:"access_token"`
}

const UserIdContextKey JWTContextKey = "user_id"

var accessTokenExpirationInMinute int = 60

func IssueAccessToken(userId uint) (string, error) {
	expirationTime := time.Now().Add(time.Minute * time.Duration(accessTokenExpirationInMinute))
	claims := &JwtClaims{
		UserId: uint(userId),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	secretKey, err := getSecretKey(os.Getenv("JWT_SECRET"))

	if err != nil {
		return "", err
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := accessToken.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	VerifyAndParseAccessToken(tokenStr)

	return tokenStr, nil
}

func VerifyAndParseAccessToken(token string) (*JwtClaimsContent, error) {
	var err error
	secretKey, err := getSecretKey(os.Getenv("JWT_SECRET"))

	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claimsContent := JwtClaimsContent{
		UserId: uint(claims["user_id"].(float64)),
		Exp:    uint64(claims["exp"].(float64)),
	}

	return &claimsContent, nil
}

func getSecretKey(key string) ([]byte, error) {
	secretKey := os.Getenv("JWT_SECRET")

	if secretKey == "" {
		return nil, errors.New("secret key is empty")
	}

	return []byte(secretKey), nil
}
