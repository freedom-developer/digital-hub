package handlers

import (
	"myapp/database"
	"myapp/models"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type MusicHandler struct {
	musicDir string
}

func NewMusicHandler(musicDir string) *MusicHandler {
	return &MusicHandler{
		musicDir: musicDir,
	}
}

// 获取音乐列表
func (h *MusicHandler) GetMusicList(c *gin.Context) {
	var musicList []models.Music

	db := database.GetDB()
	result := db.Order("id ASC").Find(&musicList)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取音乐列表失败: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    musicList,
		"message": "获取成功",
	})
}

// 播放音乐
func (h *MusicHandler) PlayMusic(c *gin.Context) {
	id := c.Param("id")

	var music models.Music
	db := database.GetDB()
	result := db.First(&music, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "音乐不存在",
		})
		return
	}

	// 返回音乐信息（实际应用中可以返回音频流）
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "播放音乐: " + music.Name,
		"data": gin.H{
			"url": music.FilePath,
		},
	})
}

// 下载音乐文件
func (h *MusicHandler) DownloadMusic(c *gin.Context) {
	id := c.Param("id")

	var music models.Music
	db := database.GetDB()
	result := db.First(&music, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "音乐不存在",
		})
		return
	}

	// 构建完整文件路径
	fullPath := filepath.Join(h.musicDir, filepath.Base(music.FilePath))

	// 发送文件
	c.File(fullPath)
}

// 获取用户信息（原有接口）
func GetUser(c *gin.Context) {
	user := models.User{
		ID:   1,
		Name: "张三",
		Age:  25,
	}

	c.JSON(http.StatusOK, user)
}
