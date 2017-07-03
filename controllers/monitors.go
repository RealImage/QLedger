package controllers

import (
	"fmt"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	response := `{"ping": "pong"}`
	fmt.Fprint(w, response)
	return
}
