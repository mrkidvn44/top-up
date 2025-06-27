package schema

type UserLoginRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required,min=10,max=10"`
	Password    string `json:"password" validate:"required,min=8"`
}

type UserCreateRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=50"`
	LastName    string `json:"last_name" validate:"required,min=2,max=50"`
	PhoneNumber string `json:"phone_number" validate:"required,min=10,max=10,numeric"`
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
