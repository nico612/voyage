package model

import (
	"time"
)

type UserM struct {
	ID        int64     `gorm:"column:id;primary_key" json:"ID,omitempty"` //
	Username  string    `gorm:"column:username" json:"username,omitempty"` //
	Password  string    `gorm:"column:password" json:"password,omitempty"` //
	Nickname  string    `gorm:"column:nickname" json:"nickname,omitempty"` //
	Email     string    `gorm:"column:email" json:"email,omitempty"`       //
	Phone     string    `gorm:"column:phone" json:"phone,omitempty"`       //
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`         //
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`         //
}

// TableName sets the insert table name for this struct type
func (u *UserM) TableName() string {
	return "user"
}
