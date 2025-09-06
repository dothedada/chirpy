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

type reqBody struct {
	Text string `json:"body"`
}

func isValidRequest(req *http.Request) (string, reqError) {
	defer req.Body.Close()

	if !strings.HasPrefix(req.Header.Get("Content-Type"), "application/json") {
		return "", reqError{
			status: http.StatusBadRequest,
			error:  errors.New("Something went wrong"),
		}
	}

	var body reqBody
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&body)
	if err != nil {
		return "", reqError{
			status: http.StatusBadRequest,
			error:  errors.New("Something went wrong"),
		}
	}

	if len(body.Text) > 140 {
		return "", reqError{
			status: http.StatusBadRequest,
			error:  errors.New("Chirp is too long"),
		}
	}

	return body.Text, reqError{}
}

func profanityCleaner(msg string, badWords []string) string {
	words := strings.Split(msg, " ")
	cleanMsg := []string{}

	for _, word := range words {
		cleanWord := word
		for _, badWord := range badWords {
			if strings.EqualFold(word, badWord) {
				cleanWord = "****"
				break
			}
		}

		cleanMsg = append(cleanMsg, cleanWord)
	}

	return strings.Join(cleanMsg, " ")
}
