package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ErrorResponse struct {
	Status int    `json:"-"`
	Code   string `json:"error_code"`
	Msg    string `json:"error_msg"`
}

func PanicInternalError(code, msg string) {
	panic(&ErrorResponse{http.StatusInternalServerError, code, msg})
}

func PanicBadRequest(code, msg string) {
	panic(&ErrorResponse{http.StatusBadRequest, code, msg})
}

type RequestHandler func() (status int, response interface{})

type Json struct {
	SetNoCacheHeaders bool
}

func (c *Json) ParseJsonBody(r *http.Request, v interface{}) {
	inData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		PanicInternalError("http_err", err.Error())
	}

	err = json.Unmarshal(inData, v)
	if err != nil {
		PanicBadRequest("json_err", err.Error())
	}
	
}

func (c *Json) Handle(w http.ResponseWriter, r *http.Request, handler RequestHandler) {
	var (
		status   int
		response interface{}
		data     []byte
		err      error
	)

	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(*ErrorResponse)
			status = err.Status
			response = err
		}

		if response != nil {
			data, err = json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		c.writeResponse(w, r, status, data)
	}()

	status, response = handler()
}

func (c *Json) writeResponse(w http.ResponseWriter, r *http.Request, status int, data []byte) {
	if c.SetNoCacheHeaders && r.Method == "GET" {
		w.Header().Set("Cache-Control", "max-age=0, no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	}
	if data != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	w.WriteHeader(status)
	if data != nil {
		w.Write(data)
	}
}
