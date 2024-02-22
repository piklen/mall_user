package dao

import (
	"fmt"
	"user/model"
)

func migration() {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.User{},
			&model.Product{},
			&model.Category{},
			&model.Address{},
			&model.Favorite{},
			&model.Notice{},
			&model.Order{},
			&model.ProductImg{},
			&model.Cart{},
			&model.Admin{},
			&model.Carousel{},
			&model.BatchUser{},
			&model.SeckillGood2MQ{},
			&model.SeckillGoods{})
	if err != nil {
		fmt.Println("err", err)
	}
	return
}
