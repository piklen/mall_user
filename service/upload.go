package service

import (
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"user/conf"
)

// UploadProductToLocalStatic 上传到本地文件中
func UploadProductToLocalStatic(file multipart.File, bossId uint, productName string) (filePath string, err error) {
	bId := strconv.Itoa(int(bossId))
	basePath := "." + conf.ProductPath + "boss" + bId + "/"
	if !DirExistOrNot(basePath) {
		CreateDir(basePath)
	}
	productPath := basePath + productName + ".jpg"
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(productPath, content, 0666)
	if err != nil {
		return "", err
	}
	return "boss" + bId + "/" + productName + ".jpg", err
}

// UploadAvatarToLocalStatic 上传头像
func UploadAvatarToLocalStatic(content []byte, userId string, userName string) (filePath string, err error) {
	basePath := "." + conf.AvatarPath + "user" + userId + "/"
	if !DirExistOrNot(basePath) {
		CreateDir(basePath)
	}
	avatarPath := basePath + userName + ".jpg"
	err = os.WriteFile(avatarPath, content, 0666)
	if err != nil {
		return "", err
	}
	return conf.PhotoHost + conf.HttpPort + conf.AvatarPath + "user" + userId + "/" + userName + ".jpg", err
}

// DirExistOrNot 判断文件夹路径是否存在
func DirExistOrNot(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// CreateDir 创建文件夹
func CreateDir(dirName string) bool {
	err := os.MkdirAll(dirName, 7550)
	if err != nil {
		return false
	}
	return true
}
