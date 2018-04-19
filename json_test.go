package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func initialize(t *testing.T, handler RequestHandler) *http.ServeMux {
	m := http.NewServeMux()
	return m
}

type id struct {
	ID string `json:"id"`
}

func TestOK(t *testing.T) {
	testCases := []struct {
		Method   string
		Code     int
		Response interface{}
	}{
		{"GET", http.StatusTeapot, nil},
		{"GET", http.StatusNotFound, id{"1234"}},
		//
		{"POST", http.StatusCreated, id{"4321"}},
	}

	for idx, tc := range testCases {
		var (
			err          error
			req          *http.Request
			jsonExpected = []byte{}
		)

		if tc.Response != nil {
			jsonExpected, err = json.Marshal(tc.Response)
			if err != nil {
				t.Fatalf("[%d] Expected nil, but got: %s", idx, err.Error())
			}
		}

		m := http.NewServeMux()

		m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			j := Json{}
			if tc.Method == "GET" {
				// Simply send the code and the response
				j.Handle(w, r, func() (status int, response interface{}) {
					return tc.Code, tc.Response
				})
			} else {
				// Parse json body and send it as response
				j.Handle(w, r, func() (status int, response interface{}) {
					var d id
					j.ParseJsonBody(r, &d)
					return tc.Code, d
				})
			}
		})

		if tc.Method == "GET" {
			req, err = http.NewRequest(tc.Method, "/", nil)
		} else {
			req, err = http.NewRequest(tc.Method, "/", bytes.NewReader(jsonExpected))
		}

		if err != nil {
			t.Fatalf("[%d] Expected nil, but got: %s", idx, err.Error())
		}

		res := httptest.NewRecorder()
		m.ServeHTTP(res, req)

		if res.Code != tc.Code {
			t.Fatalf("[%d] Expected %d, but got: %d", idx, tc.Code, res.Code)
		}

		jsonReceived, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("[%d] Expected nil, but got: %s", idx, err.Error())
		}

		if !reflect.DeepEqual(jsonExpected, jsonReceived) {
			t.Fatalf("[%d] Expected %+v, but got: %+v", idx, jsonExpected, jsonReceived)
		}
	}
}

func TestPostErrReq(t *testing.T) {
	testCases := []struct {
		Code     int
		Body     string
		Response interface{}
	}{
		{http.StatusBadRequest, `{"]`, &ErrorResponse{http.StatusBadRequest, JSONErr, "unexpected end of JSON input"}},
	}

	for idx, tc := range testCases {
		m := http.NewServeMux()

		m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			j := Json{}
			// Parse json body and send it as response
			j.Handle(w, r, func() (status int, response interface{}) {
				var d id
				j.ParseJsonBody(r, &d)
				return tc.Code, d
			})
		})

		req, err := http.NewRequest("POST", "/", strings.NewReader(tc.Body))
		if err != nil {
			t.Fatalf("[%d] Expected nil, but got: %s", idx, err.Error())
		}

		res := httptest.NewRecorder()
		m.ServeHTTP(res, req)

		if res.Code != tc.Code {
			t.Fatalf("[%d] Expected %d, but got: %d", idx, tc.Code, res.Code)
		}

		jsonExpected, err := json.Marshal(tc.Response)
		if err != nil {
			t.Fatalf("[%d] Expected nil, but got: %s", idx, err.Error())
		}

		jsonReceived, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("[%d] Expected nil, but got: %s", idx, err.Error())
		}

		if !reflect.DeepEqual(jsonExpected, jsonReceived) {
			t.Fatalf("[%d] Expected %+v, but got: %+v", idx, jsonExpected, jsonReceived)
		}
	}
}
