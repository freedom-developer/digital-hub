package user

import (
	"net/http"

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
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "用户登录成功",
		"data":    userResp,
	})
}
