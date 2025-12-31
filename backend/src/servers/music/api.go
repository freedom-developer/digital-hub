package music

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取整个音乐列表
func (ms *MusicService) GetMusicList(c *gin.Context) {
	musicList, err := ms.getMusicList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取音乐列表失败：" + err.Error(),
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
func (ms *MusicService) PlayMusic(c *gin.Context) {
	id := c.Param("id")

	music, err := ms.getMusicByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "音乐不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "播放音乐: " + music.Name,
		"data": gin.H{
			"url": music.FilePath,
		},
	})
}

// 下载音乐文件
func (ms *MusicService) DownloadMusic(c *gin.Context) {
	id := c.Param("id")
	music, err := ms.getMusicByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "音乐不存在",
		})
		return
	}

	// 构建完整文件路径
	fullPath := filepath.Join(ms.cfg.MusicDir, filepath.Base(music.FilePath))

	// 发送文件
	c.File(fullPath)
}

func (ms *MusicService) AddToFavorite(c *gin.Context) {
	userID := c.GetString("user_id")
	var req struct {
		MusicID uint `json:"music_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := ms.addToFavorite(userID, req.MusicID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "收藏成功"})
}

func (ms *MusicService) RemoveFromFavorite(c *gin.Context) {
	userID := c.GetString("user_id")
	musicID, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的音乐ID"})
		return
	}

	if err := ms.removeFromFavorite(userID, uint(musicID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "取消收藏成功"})
}

func (ms *MusicService) GetFavoriteMusic(c *gin.Context) {
	userID := c.GetString("user_id")

	musicList, err := ms.getUserMusicList(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取收藏列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    musicList,
		"message": "获取收藏列表成功",
	})
}

func (ms *MusicService) GetFavoriteMusicIDs(c *gin.Context) {
	userID := c.GetString("user_id")
	ids, err := ms.getUserMusicIDs(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取收藏列表ID失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ids})
}

func (ms *MusicService) CheckFavorite(c *gin.Context) {
	userID := c.GetString("user_id")
	musicID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的音乐ID"})
		return
	}
	isFav, err := ms.isInMyMusic(userID, uint(musicID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_favorite": isFav})
}
