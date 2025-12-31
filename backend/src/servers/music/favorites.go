package music

import (
	"errors"
	"myapp/servers/user"
	"time"

	"gorm.io/gorm"
)

type UserMusic struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"type:varchar(255);not null;index:idx_user_music_unique,unique" json:"user_id"` // 修正：UseID -> UserID
	MusicID   uint      `gorm:"not null;index:idx_user_music_unique,unique" json:"music_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// 关联关系（可选，需要时才使用）
	User  *user.User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`   // 修正：添加 references，使用指针
	Music *Music     `gorm:"foreignKey:MusicID;references:ID" json:"music,omitempty"` // 修正：使用指针
}

func (UserMusic) TableName() string {
	return "user_music"
}

// 收藏音乐
func (ms *MusicService) addToMusic(userID string, musicID uint) error {
	userMusic := &UserMusic{
		UseID:   userID,
		MusicID: musicID,
	}
	return ms.db.Create(userMusic).Error
}

// 取消收藏
func (ms *MusicService) removeFromMusic(userID string, musicID string) error {
	return ms.db.Where("user_id = ? AND music_id = ?", userID, musicID).Delete(&UserMusic{}).Error
}

// 获取用户的收藏列表
func (ms *MusicService) getUserMusicList(userID string) ([]Music, error) {
	var musicList []Music
	err := ms.db.Table("musics").
		Joins("JOIN user_music ON music.id == user_music.music_id").
		Where("user_music.user_id = ?", userID).
		Order("user_music.created_at DESC").
		Find(&musicList).Error
	return musicList, err
}

// 检查是否已收藏
func (ms *MusicService) isInMyMusic(userID string, musicID uint) (bool, error) {
	var count int64
	err := ms.db.Model(&UserMusic{}).
		Where("user_id = ? AND music_id = ?", userID, musicID).
		Count(&count).Error
	return count > 0, err
}

// 获取用户收藏的音乐ID列表（用于前端标记）
func (ms *MusicService) getUserMusicIDs(userID string) ([]uint, error) {
	var ids []uint
	err := ms.db.Model(&UserMusic{}).
		Where("user_id = ?", userID).
		Pluck("music_id", &ids).Error
	return ids, err
}

func (ms MusicService) addToFavorite(userID string, musicID uint) error {
	var music Music
	if err := ms.db.First(&music, musicID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(("音乐不存在"))
		}
		return err
	}
	fav, _ := ms.isInMyMusic(userID, musicID)
	if fav {
		return errors.New("已经收藏过该音乐")
	}

	return ms.addToMusic(userID, musicID)
}

func (ms *MusicService) removeFromFavorite(userID string, musicID uint) error {
	result := ms.db.Where("user_id = ? AND music_id = ?", userID, musicID).
		Delete(&UserMusic{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("未找到收藏记录")
	}

	return nil
}
