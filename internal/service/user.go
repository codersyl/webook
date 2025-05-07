package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"webook_Rouge/internal/domain"
	"webook_Rouge/internal/repository"
)

var (
	ErrUserDuplicatedEmail   = repository.ErrUserDuplicatedEmail
	ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码错误")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {

	encrypted, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil { // 加密出错
		return err
	}
	u.Password = string(encrypted)
	return svc.repo.Create(ctx, u) // 存起来
}

func (svc *UserService) Login(ctx context.Context, u domain.User) (domain.User, error) {
	uFromRepo, err := svc.repo.FindByEmail(ctx, u.Email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(uFromRepo.Password), []byte(u.Password))
	if err != nil {
		// 打印日志
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return uFromRepo, nil
}
