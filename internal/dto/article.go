package dto

import "github.com/StefenSutandi/sharing-vision-backend/internal/model"

type CreateArticleReq struct {
	Title    string `json:"title" validate:"required,min=20,max=200"`
	Content  string `json:"content" validate:"required,min=200"`
	Category string `json:"category" validate:"required,min=3,max=100"`
	Status   string `json:"status" validate:"required,oneof=publish draft thrash"`
}

type UpdateArticleReq struct {
	Title    string `json:"title" validate:"required,min=20,max=200"`
	Content  string `json:"content" validate:"required,min=200"`
	Category string `json:"category" validate:"required,min=3,max=100"`
	Status   string `json:"status" validate:"required,oneof=publish draft thrash"`
}

type PaginationMeta struct {
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Total  int64 `json:"total"`
}

type ListArticleRes struct {
	Data       []model.Article `json:"data"`
	Pagination PaginationMeta  `json:"pagination"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}
