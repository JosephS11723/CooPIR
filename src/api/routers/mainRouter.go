package routers

import (
	"net/http"

	"github.com/JosephS11723/CooPIR/src/api/handlers/authentication"
	"github.com/JosephS11723/CooPIR/src/api/handlers/debug"
	"github.com/JosephS11723/CooPIR/src/api/handlers/iodb"
	"github.com/JosephS11723/CooPIR/src/api/handlers/ioseaweed"
	"github.com/gin-gonic/gin"
)

func InitMainRouter() *gin.Engine {
	// intialize engine with default middleware (TODO: replace)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// set low memory limit for multipart forms (8 MiB)
	r.MaxMultipartMemory = 8 << 20

	r.LoadHTMLGlob("*index.html")

	// setup api version v1 routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// setup base path for api version v1
	v1 := r.Group("/api/v1")

	// DEBUG REQUESTS
	// debug ping challenge
	v1.GET("/ping", debug.PingPong)

	// seaweedfs file storage transfer routes
	v1.GET("/file", ioseaweed.SWGET)
	v1.POST("/file", ioseaweed.SWPOST)
	v1.DELETE("/file", ioseaweed.SWDELETE)

	// MONGO-DB
	v1.GET("/db/test", iodb.DbPingTest)
	v1.POST("/db/test", iodb.DbUploadTest)
	v1.GET("/db/test/find", iodb.DbFindTest)
	v1.POST("/db/test/find", iodb.DbUpdateTest)

	// Authentication
	v1.POST("/login", authentication.Login)
	v1.POST("/renew", authentication.RenewToken)
	v1.POST("/logout", authentication.Logout)

	// return handler router to main()
	return r
}