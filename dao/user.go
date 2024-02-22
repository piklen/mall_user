package dao

import (
	"context"
	"gorm.io/gorm"
	"user/model"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{NewDBClient(ctx)}
}

func NewUserDaoByDB(db *gorm.DB) *UserDao {
	return &UserDao{db}
}

// ExistOrNotByUserName 根据username判断是否存在该名字
func (dao *UserDao) ExistOrNotByUserName(userName string) (user *model.User, exist bool, err error) {
	var count int64
	err = dao.DB.Model(&model.User{}).Where("user_name=?", userName).Count(&count).Error
	if count == 0 {
		return user, false, err
	}
	err = dao.DB.Model(&model.User{}).Where("user_name=?", userName).First(&user).Error
	if err != nil {
		return user, false, err
	}
	return user, true, nil
}

// CreateUser 创建用户
func (dao *UserDao) CreateUser(user *model.User) error {
	return dao.DB.Model(&model.User{}).Create(&user).Error
}

// GetUserById 根据 id 获取用户
func (dao *UserDao) GetUserById(id uint) (user *model.User, err error) {
	err = dao.DB.Model(&model.User{}).Where("id=?", id).
		First(&user).Error
	return
}

// UpdateUserById 根据 id 更新用户信息
func (dao *UserDao) UpdateUserById(uId uint, user *model.User) error {
	return dao.DB.Model(&model.User{}).Where("id=?", uId).
		Updates(&user).Error
}

// BatchExistOrNotByUserNames  根据username判断是否存在该名字
func (dao *UserDao) BatchExistOrNotByUserNames(userNames []string) ([]*model.User, bool, error) {
	var users []*model.User
	var exists bool

	err := dao.DB.Model(&model.User{}).
		Where("user_name IN (?)", userNames).
		Find(&users).Error

	if err != nil {
		return nil, false, err
	}
	if len(users) != 0 {
		return users, true, nil
	}
	return users, exists, nil
}

// BatchCreateUsers 批量进行注册
func (dao *UserDao) BatchCreateUsers(users *[]model.User) error {
	if len(*users) == 0 {
		return nil
	}
	return dao.DB.Model(&model.User{}).Create(&users).Error
}
