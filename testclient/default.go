package testclient

import (
	"log"
	"os"

	"github.com/podhmo/go-webtest/tripperware"
)

var (
	defaultTripperwares tripperware.List
)

// TODO: logging interface

func init() {
	if os.Getenv("DEBUG") == "" {
		return
	}

	defaultTripperwares = tripperware.Stack(tripperware.DebugTrace())
	log.Println("builtin debug trace is activated")
}
