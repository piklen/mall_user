package model

type Notice struct {
	Model
	Text string `gorm:"type:text"`
}
