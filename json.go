package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type resError struct {
	Error string `json:"error"`
}
type resValid struct {
	Valid bool `json:"valid"`
}

func resWithErr(w http.ResponseWriter, err reqError) {
	if err.error != nil {
		log.Println(err.error.Error())
	}

	if err.status > 499 {
		log.Printf("Respond with 500+ to: %s", err.error.Error())
	}

	resJson(w, err.status, resError{
		Error: err.error.Error(),
	})
}

func resJson(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Cannot marshal response: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(data)
}
