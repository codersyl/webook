package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicatedEmail = errors.New("邮箱冲突")
	ErrUserNotFound        = gorm.ErrRecordNotFound
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
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok { // 与使用的数据库种类耦合，换一个非MySQL数据库，此段代码出错
		// 1062 唯一冲突错误，也就是表中已有相同邮箱了
		const uniqueIndeErrorNo uint16 = 1062
		if mysqlErr.Number == uniqueIndeErrorNo {
			return ErrUserDuplicatedEmail
		}
	}

	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
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
