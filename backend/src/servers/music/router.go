package music

func (ms *MusicService) RegisterRoutes() {
	api := ms.rg

	// 音乐接口
	api.GET("/music", ms.GetMusicList)
	api.GET("/music/play/:id", ms.PlayMusic)
	api.GET("/music/download/:id", ms.DownloadMusic)
}
