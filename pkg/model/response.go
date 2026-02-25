package model

type ResponseEntity[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

type ResponseEntityPagination[T any] struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    T               `json:"data"`
	Meta    *MetaPagination `json:"meta,omitempty"`
}

type MetaPagination struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"totalPage"`
	TotalData int `json:"totalData"`
}

type ResponseError[T any] struct {
	ResponseEntity[T]
	Path string `json:"path"`
}

type PaginationRequest struct {
	Page   int    `json:"page"   query:"page"   validate:"required,min=1"`
	Limit  int    `json:"limit"  query:"limit"  validate:"required,min=1,max=100"`
	Search string `json:"search" query:"search" validate:"omitempty"`
}
