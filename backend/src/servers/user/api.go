package user

import (
	"myapp/middleware"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gin-gonic/gin"
)

func (us *UserService) RegisterUser(c *gin.Context) {
	// Implementation for user registration
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的请求数据",
			"error":   err.Error(),
		})
		return
	}

	userResp, err := us.registerUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "注册用户失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "用户注册成功",
		"data":    userResp,
	})
}

func (us *UserService) Login(c *gin.Context) {
	// Implementation for user login
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的请求数据",
			"error":   err.Error(),
		})
		return
	}
	userResp, err := us.login(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "用户登录失败",
			"error":   err.Error(),
		})
		return
	}

	// 生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userResp.ID,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 72小时过期
	})

	tokenStr, err := token.SignedString([]byte(middleware.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "生成Token失败",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "用户登录成功",
		"data": gin.H{
			"user":  userResp,
			"token": tokenStr,
		},
	})
}

func (us *UserService) GetUserProfile(c *gin.Context) {
	// Implementation for getting user profile
	userID := c.GetString("user_id")
	if len(userID) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "未授权访问",
		})
		return
	}

	userResp, err := us.getUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "获取用户信息失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取用户信息成功",
		"data":    userResp,
	})
}
