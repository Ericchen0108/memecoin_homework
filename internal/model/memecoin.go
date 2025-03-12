package model

import (
	"time"

	"gorm.io/gorm"
)

type Memecoin struct {
	ID              int            `json:"id" gorm:"primaryKey"`
	Name            string         `json:"name" gorm:"unique; not null"`
	Description     string         `json:"description"`
	CreatedAt       time.Time      `json:"created_at"`
	PopularityScore int            `json:"popularity_score"`
	Deleted         gorm.DeletedAt `gorm:"index"`
}
