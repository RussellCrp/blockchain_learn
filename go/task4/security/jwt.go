package security

import (
	"blogs_learn/utils"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const jwtKey = "abcdefg"

type UserClaims struct {
	UserID uint
	jwt.RegisteredClaims
}

func ValidUser(c *gin.Context) (uint, bool) {
	tokenString, err := utils.GetToken(c)
	if err != nil {
		return 0, false
	}
	// 解析 token
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})
	if err != nil || !token.Valid {
		return 0, false
	}
	if claims.UserID == 0 {
		return 0, false
	}
	return claims.UserID, true
}

func GenerateToken(userID uint) (string, error) {
	// 创建 JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
