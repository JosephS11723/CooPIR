package main

import (
	"github.com/JosephS11723/CooPIR/src/api/routers"
	"log"
)

// API Documentation
// @title CooPIR API
// @version 1.0.0
// @host localhost:8080
// @BasePath /api/v1

func main() {
	// set log to print line numbers
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// initialize router with handlers
	r := routers.InitRouter()

	// run and serve
	r.Run("0.0.0.0:8080")
}
