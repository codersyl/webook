package repository

import (
	"context"
	"webook_Rouge/internal/domain"
	"webook_Rouge/internal/repository/dao"
)

var (
	ErrUserDuplicatedEmail = dao.ErrUserDuplicatedEmail
	// ErrUserDuplicatedEmailV1 = fmt.Errorf("邮箱冲突", dao.ErrUserDuplicatedEmail) // 更好的定义
	ErrUserNotFound = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}
func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Email:    u.Email,
		Password: u.Password,
	}, err
}
