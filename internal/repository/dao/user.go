package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	// 如果有缓存，则在这里操作缓存

	now := time.Now().UnixMilli()
	u.UpdatedTime = now
	u.CreatedTime = now
	return dao.db.WithContext(ctx).Create(&u).Error
}

// 直接对应数据库中的表结构
type User struct {
	Id       int64  `gorm:"primary_key,auto_increment"`
	Email    string `gorm:"unique"`
	Password string

	// 创建时间，毫秒数，后端的时间相关数据用UTC+0来存储
	CreatedTime int64
	// 更新时间，毫秒数
	UpdatedTime int64
}
