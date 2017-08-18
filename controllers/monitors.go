package controllers

import (
	"net/http"
)

// Ping responds 200 OK when the server is up and healthy
func Ping(w http.ResponseWriter, r *http.Request) {
	// TODO: Should DB connection check be made while ping ?
	response := `{"ping": "pong"}`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(response))
	return
}
