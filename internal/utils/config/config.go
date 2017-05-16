package config

import "os"

var PORT string

func init() {
	PORT = os.Getenv("PORT")
}
