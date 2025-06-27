package mapper

import (
	"top-up-api/internal/model"
	"top-up-api/internal/schema"
)

func UserProfileResponseFromModel(user *model.User) *schema.UserProfileResponse {
	return &schema.UserProfileResponse{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Balance:     user.Balance,
		CashBack:    user.CashBack,
	}
}

func UserLoginDetailFromModel(user *model.User) *schema.UserLoginDetail {
	return &schema.UserLoginDetail{
		ID:          user.ID,
		PhoneNumber: user.PhoneNumber,
		Password:    user.Password,
	}
}
func UserCreateRequestToModel(user *schema.UserCreateRequest) *model.User {
	return &model.User{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Password:    user.Password,
	}
}
