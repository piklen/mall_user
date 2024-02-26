package model

// Admin 管理员
type Admin struct {
	Model
	UserName       string
	PasswordDigest string
	Avatar         string
}
