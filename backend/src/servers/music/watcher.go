package music

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"gorm.io/gorm"
)

type FileWatcher struct {
	watcher  *fsnotify.Watcher
	musicDir string
	db       *gorm.DB
}

func NewFileWatcher(musicDir string, db *gorm.DB) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// 确保目录存在
	if err := os.MkdirAll(musicDir, 0755); err != nil {
		return nil, err
	}

	return &FileWatcher{
		watcher:  watcher,
		musicDir: musicDir,
		db:       db,
	}, nil
}

func (fw *FileWatcher) Start() error {
	// 添加监控目录
	err := fw.watcher.Add(fw.musicDir)
	if err != nil {
		return err
	}

	log.Printf("开始监控音乐目录: %s", fw.musicDir)

	// 初始化：扫描现有文件
	fw.scanExistingFiles()

	// 启动监控协程
	go fw.watch()

	return nil
}

func (fw *FileWatcher) watch() {
	for {
		select {
		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}

			log.Printf("检测到文件事件: %s - %s", event.Op, event.Name)

			// 只处理音频文件
			if !isMusicFile(event.Name) {
				log.Printf("忽略非音乐文件: %s", event.Name)
				continue
			}

			switch {
			case event.Op&fsnotify.Create == fsnotify.Create:
				fw.handleCreate(event.Name)
			case event.Op&fsnotify.Remove == fsnotify.Remove:
				fw.handleDelete(event.Name)
			case event.Op&fsnotify.Rename == fsnotify.Rename:
				fw.handleDelete(event.Name)
			}

		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("文件监控错误: %v", err)
		}
	}
}

func (fw *FileWatcher) scanExistingFiles() {
	files, err := os.ReadDir(fw.musicDir)
	if err != nil {
		log.Printf("扫描目录失败: %v", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(fw.musicDir, file.Name())
		if isMusicFile(filePath) {
			fw.handleCreate(filePath)
		}
	}

	log.Printf("已扫描现有音乐文件，共 %d 个", len(files))
}

func (fw *FileWatcher) handleCreate(filePath string) {
	fileName := getFileName(filePath)
	relativePath := getRelativePath(fw.musicDir, filePath)

	// 检查是否已存在
	var existing Music
	result := fw.db.Where("name = ?", fileName).First(&existing)
	if result.Error == nil {
		log.Printf("音乐已存在: %s", fileName)
		return
	}

	// 添加到数据库
	music := Music{
		Name:     fileName,
		FilePath: relativePath,
	}

	if err := fw.db.Create(&music).Error; err != nil {
		log.Printf("添加音乐失败: %v", err)
		return
	}

	log.Printf("✅ 新增音乐: %s (ID: %d)", fileName, music.ID)
}

func (fw *FileWatcher) handleDelete(filePath string) {
	fileName := getFileName(filePath)

	result := fw.db.Where("name = ?", fileName).Delete(&Music{})
	if result.Error != nil {
		log.Printf("删除音乐失败: %v", result.Error)
		return
	}

	if result.RowsAffected > 0 {
		log.Printf("🗑️  删除音乐: %s", fileName)
	}
}

func (fw *FileWatcher) handleRename(oldFile, newFile string) {
	// newName := getFileName(newFile)
	// oldName := getFileName(oldFile)

	// db := database.GetDB()

	// result := db.Where("name = ?", oldName).Update(&models.Music{.Name = newName})
	// if result.Error != nil {
	// 	log.Printf("删除音乐失败: %v", result.Error)
	// 	return
	// }

	// if result.RowsAffected > 0 {
	// 	log.Printf("🗑️  删除音乐: %s", fileName)
	// }
}

func (fw *FileWatcher) Close() error {
	return fw.watcher.Close()
}

// 辅助函数
func isMusicFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	musicExts := []string{".mp3", ".flac", ".wav", ".aac", ".ogg", ".m4a"}
	for _, musicExt := range musicExts {
		if ext == musicExt {
			return true
		}
	}
	return false
}

func getFileName(filePath string) string {
	base := filepath.Base(filePath)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

func getRelativePath(baseDir, filePath string) string {
	rel, err := filepath.Rel(baseDir, filePath)
	if err != nil {
		return filePath
	}
	return "/" + filepath.ToSlash(rel)
}
