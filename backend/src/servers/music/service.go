package music

import (
	"context"
	logger "myapp/log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MusicConfig struct {
	MusicDir string
}

type MusicService struct {
	cfg         *MusicConfig
	db          *gorm.DB
	fileWatcher *FileWatcher
	rg          *gin.RouterGroup
}

func NewMusicService(ctx context.Context, cfg *MusicConfig, db *gorm.DB, rg *gin.RouterGroup) *MusicService {
	watcher, err := NewFileWatcher(cfg.MusicDir, db)
	if watcher == nil {
		logger.ZError(&ctx, "创建文件监控器失败", err)
		return nil
	}

	// 自动迁移
	err = db.AutoMigrate(&Music{})
	if err != nil {
		// logger.ZError(&ctx, "数据库自动迁移失败", err)
		return nil
	}

	return &MusicService{
		cfg:         cfg,
		db:          db,
		fileWatcher: watcher,
		rg:          rg,
	}
}

func (ms *MusicService) Start() {
	// 启动音乐服务的逻辑
	if ms.fileWatcher != nil {
		if err := ms.fileWatcher.Start(); err != nil {
			logger.ZError(nil, "启动文件监控失败", err)
		}
	} else {
		logger.ZError(nil, "文件监控器未初始化", nil)
	}

	ms.RegisterRoutes()
}
