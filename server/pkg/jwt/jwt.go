package jwt

import (
	"time"

	"server/infra/config"

	jwtv4 "github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwtv4.RegisteredClaims
}

func GenerateToken(id int64, username string, isAdmin bool) (string, error) {
	cfg := config.GetConfig()
	claims := Claims{
		ID:       id,
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwtv4.RegisteredClaims{
			ExpiresAt: jwtv4.NewNumericDate(time.Now().Add(time.Duration(cfg.ExpireDuration) * time.Hour)),
			Issuer:    cfg.Issuer,
			Subject:   cfg.Subject,
			IssuedAt:  jwtv4.NewNumericDate(time.Now()),
		},
	}

	token := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Key))
}

func ParseToken(token string) (*Claims, bool) {
	claims := new(Claims)
	parsedToken, err := jwtv4.ParseWithClaims(token, claims, func(token *jwtv4.Token) (interface{}, error) {
		return []byte(config.GetConfig().Key), nil
	})
	if err != nil || parsedToken == nil || !parsedToken.Valid || claims == nil {
		return nil, false
	}
	return claims, true
}
