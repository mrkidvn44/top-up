package mapper

import "top-up-api/internal/schema"

func PaginationResponseFromModel(totalCount, totalPage, currentPage int, data interface{}) *schema.PaginationResponse {
	return &schema.PaginationResponse{
		Code:       200,
		Message:    "success",
		Data:       data,
		Pagination: schema.Pagination{TotalCount: totalCount, TotalPage: totalPage, CurrentPage: currentPage},
	}
}

func SuccessResponse(data interface{}) *schema.Response {
	return &schema.Response{Code: 200, Message: "success", Data: data}
}

func ErrorResponse(code int, message string, error string) *schema.Response {
	return &schema.Response{Code: code, Message: message, Error: error}
}
