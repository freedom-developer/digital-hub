package music

import (
	"time"
)

type Music struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null;index;unique" json:"name"`
	FilePath  string    `gorm:"type:varchar(1024);not null" json:"file_path"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (ms *MusicService) getMusicList() ([]Music, error) {
	var musicList []Music
	result := ms.db.Order("id ASC").Find(&musicList)
	if result.Error != nil {
		return nil, result.Error
	}
	return musicList, nil
}

func (ms *MusicService) getMusicByID(id string) (*Music, error) {
	var music Music
	result := ms.db.First(&music, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &music, nil
}
