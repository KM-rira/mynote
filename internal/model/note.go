package model

import "time"

type Note struct {
	ID        int       `gorm:"primaryKey"`
	Title     string    `json:"title"`
	Contents  string    `json:"contents"`
	Category  string    `json:"category"`
	Important bool      `json:"important"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (Note) TableName() string {
	return "note"
}
