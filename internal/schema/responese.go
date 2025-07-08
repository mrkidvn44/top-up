package schema

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}


type PaginationResponse struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Error      string      `json:"error"`
	Pagination Pagination  `json:"pagination"`
	Data       interface{} `json:"data"`
}

type Pagination struct {
	TotalCount  int `json:"total_count"`
	TotalPage   int `json:"total_page"`
	CurrentPage int `json:"current_page"`
}

