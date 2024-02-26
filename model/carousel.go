package model

// Carousel 轮播图
type Carousel struct {
	Model
	ImgPath   string
	ProductId uint `gorm:"not null"`
}
