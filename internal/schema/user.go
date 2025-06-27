package schema

import "top-up-api/internal/model"

type UserLoginRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required,min=10,max=10"`
	Password    string `json:"password" validate:"required,min=8"`
}

type UserCreateRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=50"`
	LastName    string `json:"last_name" validate:"required,min=2,max=50"`
	PhoneNumber string `json:"phone_number" validate:"required,min=10,max=10"`
	Password    string `json:"password" validate:"required,min=8"`
}

type UserLoginDetail struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type UserAuthDetail struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phone_number"`
}

type UserProfileResponse struct {
	ID          uint   `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Balance     int    `json:"balance"`
	CashBack    int    `json:"cash_back"`
}

func UserProfileResponseFromModel(user *model.User) *UserProfileResponse {
	return &UserProfileResponse{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Balance:     user.Balance,
		CashBack:    user.CashBack,
	}
}

func UserLoginDetailFromModel(user *model.User) *UserLoginDetail {
	return &UserLoginDetail{
		ID:          user.ID,
		PhoneNumber: user.PhoneNumber,
		Password:    user.Password,
	}
}
func UserCreateRequestToModel(user *UserCreateRequest) *model.User {
	return &model.User{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Password:    user.Password,
	}
}
