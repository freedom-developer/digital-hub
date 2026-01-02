package middleware

import (
	"log"
	"os"
)

var JWTSecret string

func init() {
	// 从环境变量获取 JWT 密钥
	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		// 开发环境默认值（生产环境必须设置环境变量）
		JWTSecret = "jwt_default_secret"
		log.Println("警告: 使用默认 JWT 密钥，请在生产环境中更改")
	}
}
