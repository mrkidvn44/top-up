package service

import (
	"context"
	"errors"
	"top-up-api/internal/mapper"
	"top-up-api/internal/repository"
	"top-up-api/internal/schema"
	"top-up-api/pkg/auth"

	"github.com/jackc/pgx/v5/pgconn"
)

type IUserService interface {
	GetUserByID(ctx context.Context, id uint) (*schema.UserProfileResponse, error)
	Login(ctx context.Context, phoneNumber string, password string) (*schema.UserLoginDetail, error)
	CreateUser(ctx context.Context, user schema.UserCreateRequest) error
}
type UserService struct {
	repo repository.IUserRepository
}

var _ IUserService = (*UserService)(nil)

func NewUserService(userRepository repository.IUserRepository) *UserService {
	return &UserService{repo: userRepository}
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*schema.UserProfileResponse, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.UserProfileResponseFromModel(user), nil
}

func (s *UserService) Login(ctx context.Context, phoneNumber string, password string) (*schema.UserLoginDetail, error) {
	user, err := s.repo.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, err
	}
	if !auth.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("wrong password")
	}
	return mapper.UserLoginDetailFromModel(user), nil
}

func (s *UserService) CreateUser(ctx context.Context, user schema.UserCreateRequest) error {
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	err = s.repo.CreateUser(ctx, mapper.UserCreateRequestToModel(&user))

	var perr *pgconn.PgError
	errors.As(err, &perr)
	if err != nil {
		if perr.Code == "23505" {
			return errors.New("user already exists")
		}
		return err
	}
	return nil
}
