package music

import "myapp/middleware"

func (ms *MusicService) RegisterRoutes() {
	musicGroup := ms.rg

	// 音乐接口
	musicGroup.GET("", ms.GetMusicList)
	musicGroup.GET("/play/:id", ms.PlayMusic)
	musicGroup.GET("/download/:id", ms.DownloadMusic)

	// 收藏
	favGroup := musicGroup.Group("/favorite")
	favGroup.Use(middleware.AuthMiddleware())
	favGroup.POST("", ms.AddToFavorite)            // 添加收藏
	favGroup.DELETE("/:id", ms.RemoveFromFavorite) // 取消收藏
	favGroup.GET("", ms.GetFavoriteMusic)          // 获取收藏列表
	favGroup.GET("/ids", ms.GetFavoriteMusicIDs)   // 获取收藏ID列表
	favGroup.GET("/check/:id", ms.CheckFavorite)   // 检查是否收藏
}
