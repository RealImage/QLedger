package controllers

import (
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	response := `{"ping": "pong"}`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(response))
	return
}
