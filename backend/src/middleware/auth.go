// middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未提供认证令牌",
			})
			c.Abort()
			return
		}

		// 2. 检查格式是否为 "Bearer {token}"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "认证令牌格式错误",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. 解析 JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的认证令牌",
			})
			c.Abort()
			return
		}

		// 4. 验证 token 是否有效
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "认证令牌已失效",
			})
			c.Abort()
			return
		}

		// 5. 提取 claims 中的 user_id
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID, exists := claims["user_id"]
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "令牌中缺少用户信息",
				})
				c.Abort()
				return
			}

			// 6. 将 user_id 存入 gin 上下文（重要！）
			c.Set("user_id", userID.(string))

			// 可选：打印日志
			// log.Printf("用户认证成功 - UserID: %s", userID)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无法解析令牌",
			})
			c.Abort()
			return
		}

		// 7. 继续处理请求
		c.Next()
	}
}
