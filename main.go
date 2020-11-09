package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mikhail-angelov/telegram-bot-blur/web"
)

func getEnv(key string, def string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return val
}

func main() {
	p := getEnv("HTTP_PORT", "8008")
	s := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", p),
		Handler: api.GetRouter(),
	}
	log.Printf("Listening on %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}
