package dao

import (
	"context"
	"gorm.io/gorm"
	"user/model"
)

type BatchUserDao struct {
	*gorm.DB
}

func BatchNewUserDao(ctx context.Context) *BatchUserDao {
	return &BatchUserDao{NewDBClient(ctx)}
}

func BatchNewUserDaoByDB(db *gorm.DB) *BatchUserDao {
	return &BatchUserDao{db}
}

// ExistOrNotByUserNames 根据username判断是否存在该名字
func (dao *BatchUserDao) BatchExistOrNotByUserNames(userNames []string) ([]*model.BatchUser, bool, error) {
	var users []*model.BatchUser
	var exists bool

	err := dao.DB.Model(&model.User{}).
		Where("user_name IN (?)", userNames).
		Find(&users).Error

	if err != nil {
		return nil, false, err
	}

	for _, user := range users {
		if user == nil {

		} else {
			exists = true
		}
	}

	return users, exists, nil
}

// CreateUsers 批量进行注册
func (dao *BatchUserDao) BatchCreateUsers(users *[]model.BatchUser) error {
	if len(*users) == 0 {
		return nil
	}
	return dao.DB.Model(&model.User{}).Create(&users).Error
	//return dao.DB.Create(&users).Error
}
