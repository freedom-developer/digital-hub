package main

import (
	"log"
	"myapp/config"
	"myapp/database"
	"myapp/handlers"
	"myapp/watcher"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 启动文件监控
	fileWatcher, err := watcher.NewFileWatcher(cfg.MusicDir)
	if err != nil {
		log.Fatalf("创建文件监控器失败: %v", err)
	}

	if err := fileWatcher.Start(); err != nil {
		log.Fatalf("启动文件监控失败: %v", err)
	}
	defer fileWatcher.Close()

	// 创建 Gin 实例
	r := gin.Default()

	// 配置 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 创建处理器
	musicHandler := handlers.NewMusicHandler(cfg.MusicDir)

	// 路由
	api := r.Group("/api")
	{
		// 用户接口
		api.GET("/user", handlers.GetUser)

		// 音乐接口
		api.GET("/music", musicHandler.GetMusicList)
		api.GET("/music/play/:id", musicHandler.PlayMusic)
		api.GET("/music/download/:id", musicHandler.DownloadMusic)
	}

	// 优雅关闭
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("收到关闭信号，正在关闭服务...")
		fileWatcher.Close()
		os.Exit(0)
	}()

	// 启动服务
	addr := "0.0.0.0:" + cfg.ServerPort
	log.Printf("🚀 服务器启动在: http://%s", addr)
	log.Printf("📁 监控音乐目录: %s", cfg.MusicDir)

	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
