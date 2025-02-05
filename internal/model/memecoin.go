package model

import "time"

type Memecoin struct {
	ID              int        `json:"id" gorm:"primaryKey"`
	Name            string     `json:"name" gorm:"unique; not null"`
	Description     string     `json:"description"`
	CreatedAt       time.Time  `json:"created_at"`
	PopularityScore int        `json:"popularity_score"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}
