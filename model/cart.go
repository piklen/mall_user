package model

// Cart 购物车
type Cart struct {
	Model
	UserId    uint `gorm:"not null"`
	ProductId uint `gorm:"not null"`
	BossId    uint `gorm:"not null"`
	Num       uint `gorm:"not null"`
	MaxNum    uint `gorm:"not null"`
	Check     bool
}
