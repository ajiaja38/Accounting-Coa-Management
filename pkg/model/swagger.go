package model

type SwaggerAuthResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type SwaggerCOAResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type SwaggerCOAListResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    any             `json:"data"`
	Meta    *MetaPagination `json:"meta,omitempty"`
}

type SwaggerJournalResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type SwaggerJournalListResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    any             `json:"data"`
	Meta    *MetaPagination `json:"meta,omitempty"`
}

type SwaggerEmptyResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SwaggerErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Path    string `json:"path"`
}
