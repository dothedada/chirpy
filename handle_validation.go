package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type reqError struct {
	status int
	error  error
}

func isValidRequest(req *http.Request) reqError {
	if !strings.HasPrefix(req.Header.Get("Content-Type"), "application/json") {
		return reqError{
			status: http.StatusBadRequest,
			error:  errors.New("Something went wrong"),
		}
	}

	type body struct {
		Text string `json:"body"`
	}
	var reqBody body
	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		return reqError{
			status: http.StatusBadRequest,
			error:  errors.New("Something went wrong"),
		}
	}

	if len(reqBody.Text) > 140 {
		return reqError{
			status: http.StatusBadRequest,
			error:  errors.New("Chirp is too long"),
		}
	}

	return reqError{}
}
