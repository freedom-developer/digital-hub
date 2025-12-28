package music

import "time"

type Music struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null;index;unique" json:"name"`
	FilePath  string    `gorm:"type:varchar(1024);not null" json:"file_path"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
