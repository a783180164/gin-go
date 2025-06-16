// pkg/jwt/jwt.go
package jwtmidd

import (
	"net/http"
	"strings"
	"time"

	"gin-go/configs"
	"gin-go/pkg/code"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var config = configs.Get()

// JWTConfig holds configuration for JWT generation and validation.
type JWTConfig struct {
	SecretKey     string        // 用于签名和校验的密钥
	ExpiresIn     time.Duration // Token 过期时间
	AuthHeaderKey string        // 请求头中存放 Token 的字段
}

// DefaultJWTConfig 返回一个默认配置
func DefaultJWTConfig() *JWTConfig {
	return &JWTConfig{
		SecretKey:     config.JWT.Secret,
		ExpiresIn:     time.Hour * time.Duration(config.JWT.Hour),
		AuthHeaderKey: "Authorization",
	}
}

// GenerateToken 根据用户 ID 或其他 claim 生成一个 JWT token。
func (cfg *JWTConfig) GenerateToken(claims jwt.MapClaims) (string, error) {
	claims["exp"] = time.Now().Add(cfg.ExpiresIn).Unix()
	// 使用 HMAC SHA256 签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.SecretKey))
}

// ValidateToken 解析并校验 token 字符串，返回 token 对象和 claims。
func (cfg *JWTConfig) ValidateToken(tokenStr string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(cfg.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, jwt.ErrInvalidKey
	}
	return token, claims, nil
}

// AuthMiddleware 返回一个 Gin 中间件，用于验证 JWT。
func (cfg *JWTConfig) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 拿 token: "Bearer <token>"
		authHeader := c.GetHeader(cfg.AuthHeaderKey)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusOK, &code.Failure{
				Code:    code.AuthorizationNo,
				Data:    nil,
				Message: code.Text(code.AuthorizationNo),
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusOK, &code.Failure{
				Code:    code.AuthorizationError,
				Data:    nil,
				Message: code.Text(code.AuthorizationError),
			})
			return
		}

		tokenStr := parts[1]
		_, claims, err := cfg.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, &code.Failure{
				Code:    code.AuthorizationError,
				Data:    nil,
				Message: code.Text(code.AuthorizationError),
			})
			return
		}

		// 将解析后的 claims 存入 context，供后续处理使用
		c.Set("jwt_claims", claims)
		c.Next()
	}
}
