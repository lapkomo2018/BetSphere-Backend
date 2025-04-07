package model

import "time"

type User struct {
	ID        uint64    `gorm:"primary_key" json:"id"`
	Email     string    `gorm:"unique" json:"email"`
	Username  string    `gorm:"unique" json:"username"`
	Password  string    `gorm:"" json:"-"`
	Admin     bool      `gorm:"default:false" json:"admin"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
