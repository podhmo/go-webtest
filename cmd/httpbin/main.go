package main

import (
	"log"
	"os"

	"github.com/podhmo/go-webtest/httpbin"
)

func main() {
	url := os.Args[1]
	log.Fatalf("!%+v", httpbin.Run(url, httpbin.Handler()))
}
