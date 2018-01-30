package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MainfluxLabs/rules-engine/api"
)

const (
	defPort string = "9000"
	envPort string = "PORT"
)

type config struct {
	Port string
}

func main() {
	cfg := config{
		Port: getenv(envPort, defPort),
	}

	p := fmt.Sprintf(":%s", cfg.Port)
	http.ListenAndServe(p, api.MakeHandler())
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
