package service

import (
	"context"
	"errors"
	"top-up-api/internal/repository"
	"top-up-api/internal/schema"
	"top-up-api/pkg/auth"

	"github.com/jackc/pgx/v5/pgconn"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{repo: userRepository}
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*schema.UserProfileResponse, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return schema.UserProfileResponseFromModel(user), nil
}

func (s *UserService) Login(ctx context.Context, phoneNumber string, password string) (*schema.UserLoginDetail, error) {
	user, err := s.repo.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, err
	}
	if !auth.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid password")
	}
	return schema.UserLoginDetailFromModel(user), nil
}

func (s *UserService) CreateUser(ctx context.Context, user schema.UserCreateRequest) error {
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	err = s.repo.CreateUser(ctx, schema.UserCreateRequestToModel(&user))

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

func (s *UserService) DeductBalance(ctx context.Context, userID uint, amount int) error {
	return s.repo.DeductBalance(ctx, userID, amount)
}

func (s *UserService) AddBalance(ctx context.Context, userID uint, amount int) error {
	return s.repo.AddBalance(ctx, userID, amount)
}

func (s *UserService) AddCashBack(ctx context.Context, userID uint, amount int) error {
	return s.repo.AddCashBack(ctx, userID, amount)
}
