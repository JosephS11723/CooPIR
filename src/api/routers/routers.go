package routers

import (
	"github.com/JosephS11723/CooPIR/src/api/handlers/debug"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// intialize engine with default middleware (TODO: replace)
	r := gin.Default()

	// debug ping challenge
	r.GET("/ping", debug.PingPong)

	// return handler router to main()
	return r
}
