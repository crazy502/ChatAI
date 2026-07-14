package jwt

import (
	"fmt"
	"time"

	"server/infra/config"

	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Claims struct {
	ID                     int64  `json:"id"`       // 用户ID
	Username               string `json:"username"` // 用户名
	IsAdmin                bool   `json:"is_admin"` // 是否为管理员
	jwtv4.RegisteredClaims        // 注册的JWT声明
}

// GenerateToken 生成JWT令牌
func GenerateToken(id int64, username string, isAdmin bool) (string, error) {
	//1. 生成JWT声明
	cfg := config.GetConfig()
	now := time.Now().UTC()
	claims := Claims{
		// 注册的JWT声明
		ID:       id,
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwtv4.RegisteredClaims{
			// 过期时间
			ExpiresAt: jwtv4.NewNumericDate(now.Add(time.Duration(cfg.ExpireDuration) * time.Hour)),
			// 签发者
			Issuer: cfg.Issuer,
			// 主题
			Subject: cfg.Subject,
			// 签发时间
			IssuedAt:  jwtv4.NewNumericDate(now),
			NotBefore: jwtv4.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	}

	token := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Key))
}

// ParseToken 解析JWT令牌
func ParseToken(token string) (*Claims, bool) {
	//1. 解析JWT令牌
	claims := new(Claims)
	//2. 验证JWT令牌
	cfg := config.GetConfig()
	parser := jwtv4.NewParser(jwtv4.WithValidMethods([]string{jwtv4.SigningMethodHS256.Alg()}))
	parsedToken, err := parser.ParseWithClaims(token, claims, func(token *jwtv4.Token) (interface{}, error) {
		if token.Method != jwtv4.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}
		return []byte(cfg.Key), nil
	})
	//3. 检查JWT令牌是否有效
	if err != nil || parsedToken == nil || !parsedToken.Valid || claims == nil {
		return nil, false
	}
	if claims.ID <= 0 || claims.Username == "" || !claims.VerifyIssuer(cfg.Issuer, true) || claims.Subject != cfg.Subject {
		return nil, false
	}
	return claims, true
}
