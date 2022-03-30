package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/JosephS11723/CooPIR/src/api/lib/crypto"
	"github.com/JosephS11723/CooPIR/src/api/routers"
)

// API Documentation
// @title     CooPIR API
// @version   1.0.0
// @host      localhost:8080
// @BasePath  /api/v1

func main() {
	// set log to print line numbers
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// initialize router with handlers
	r := routers.InitMainRouter()

	// key check. this initializes the public and private keys for authentication (JWT token signing)
	crypto.VerifyKeys()

	// run and serve main router
	go func() {
		r.Run("0.0.0.0:8080")
	}()

	// wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 10)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
}
