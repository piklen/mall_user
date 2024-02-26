package model

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

//func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
//	u.CreatedAt = time.Now() // 设置创建时间为当前系统时间
//	return
//}
//func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
//	u.UpdatedAt = time.Now() // 设置更新时间为当前系统时间
//	return
//}
