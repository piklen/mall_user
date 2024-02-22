package dao

import (
	"context"
	"gorm.io/gorm"
	"user/model"
	"user/pkg/e"
)

type CartDao struct {
	*gorm.DB
}

func NewCartDao(ctx context.Context) *CartDao {
	return &CartDao{NewDBClient(ctx)}
}

func NewCartDaoByDB(db *gorm.DB) *CartDao {
	return &CartDao{db}
}

// CreateCart 创建 cart pId(商品 id)、uId(用户id)、bId(店家id)
func (dao *CartDao) CreateCart(pId, uId, bId uint) (cart *model.Cart, status int, err error) {
	// 查询有无此条商品
	cart, err = dao.GetCartById(pId, uId, bId)
	// 空的，第一次加入
	if err == gorm.ErrRecordNotFound {
		cart = &model.Cart{
			UserId:    uId,
			ProductId: pId,
			BossId:    bId,
			Num:       1,
			MaxNum:    10,
			Check:     false,
		}
		err = dao.DB.Create(&cart).Error
		if err != nil {
			return
		}
		return cart, e.Success, err
	} else if cart.Num < cart.MaxNum {
		// 小于最大 num
		cart.Num++
		err = dao.DB.Save(&cart).Error
		if err != nil {
			return
		}
		return cart, e.ErrorProductExistCart, err
	} else {
		// 大于最大num
		return cart, e.ErrorProductMoreCart, err
	}
}

// GetCartById 获取 Cart 通过 id
func (dao *CartDao) GetCartById(pId, uId, bId uint) (cart *model.Cart, err error) {
	err = dao.DB.Model(&model.Cart{}).
		Where("user_id=? AND product_id=? AND boss_id=?", uId, pId, bId).
		First(&cart).Error
	return
}

// ListCartByUserId 获取 Cart 通过 user_id
func (dao *CartDao) ListCartByUserId(uId uint) (cart []*model.Cart, err error) {
	err = dao.DB.Model(&model.Cart{}).
		Where("user_id=?", uId).Find(&cart).Error
	return
}

// UpdateCartNumById 通过id更新Cart信息
func (dao *CartDao) UpdateCartNumById(cId, num uint) error {
	return dao.DB.Model(&model.Cart{}).
		Where("id=?", cId).Update("num", num).Error
}

// DeleteCartById 通过 cart_id 删除 cart
func (dao *CartDao) DeleteCartById(cId uint) error {
	return dao.DB.Model(&model.Cart{}).
		Where("id=?", cId).Delete(&model.Cart{}).Error
}
