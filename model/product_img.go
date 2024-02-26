package model

type ProductImg struct {
	Model
	ProductID uint `gorm:"not null"`
	ImgPath   string
}
