package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTSecret JWT密钥（建议从配置文件读取）
var JWTSecret = []byte("your-secret-key")

// JWTClaims JWT载荷
type JWTClaims struct {
	UID      uint64 `json:"uid"`
	AgentID  uint32 `json:"agent_id"`
	Username string `json:"username"`

	jwt.RegisteredClaims
}

// GenerateToken 生成JWT
func GenerateToken(uid uint64, agentID uint32, username string) (string, error) {

	claims := JWTClaims{
		UID:      uid,
		AgentID:  agentID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(JWTSecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*JWTClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return JWTSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}