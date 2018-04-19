package controller

import "net/http"

type ErrorCode string

const (
	DBErr   ErrorCode = "db_err"
	JSONErr           = "json_err"
	HTTPErr           = "http_err"
	GenErr            = "gen_err"
)

type ErrorResponse struct {
	Status int       `json:"-"`
	Code   ErrorCode `json:"error_code"`
	Msg    string    `json:"error_msg"`
}

func PanicInternalError(code ErrorCode, msg string) {
	panic(&ErrorResponse{http.StatusInternalServerError, code, msg})
}

func PanicBadRequest(code ErrorCode, msg string) {
	panic(&ErrorResponse{http.StatusBadRequest, code, msg})
}

func PanicNotFound(code ErrorCode, msg string) {
	panic(&ErrorResponse{http.StatusNotFound, code, msg})
}
