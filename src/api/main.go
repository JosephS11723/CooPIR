package main

import (
	"github.com/JosephS11723/CooPIR/src/api/routers"
)

func main() {
	// initialize router with handlers
	r := routers.InitRouter()

	// run and server
	r.Run("localhost:8080")
}
