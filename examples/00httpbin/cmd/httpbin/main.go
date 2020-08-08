package main

import (
	"log"
	"os"

	"m/httpbin"
)

func main() {
	url := os.Args[1]
	log.Fatalf("!%+v", httpbin.Run(url, httpbin.Handler()))
}
