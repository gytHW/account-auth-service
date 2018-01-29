package model

import (
	"time"
)

// Model 模型基本 struct
type Model struct {
	//ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"CreatedAt"`
	UpdatedAt time.Time  `json:"UpdatedAt"`
	DeletedAt *time.Time `json:"DeletedAt"`
}

type User_auth struct {
	Id         int    `gorm:"primary_key"`
	Password   string `gorm:"type:varchar(100);not null"`
	Is_initial int    `gorm:"not null"`
	Model
}
