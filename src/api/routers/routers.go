package routers

import (
	"github.com/JosephS11723/CooPIR/src/api/handlers/authentication"
	"github.com/JosephS11723/CooPIR/src/api/handlers/debug"
	"github.com/JosephS11723/CooPIR/src/api/handlers/iodb"
	"github.com/JosephS11723/CooPIR/src/api/handlers/ioseaweed"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// intialize engine with default middleware (TODO: replace)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// set low memory limit for multipart forms (8 MiB)
	r.MaxMultipartMemory = 8 << 20

	// DEBUG REQUESTS
	// debug ping challenge
	r.GET("/ping", debug.PingPong)

	// seaweedfs file storage transfer routes
	r.GET("/file", ioseaweed.SWGET)
	r.POST("/file", ioseaweed.SWPOST)
	r.DELETE("/file", ioseaweed.SWDELETE)

	// MONGO-DB
	r.GET("/db/test", iodb.DbPingTest)
	r.POST("/db/test", iodb.DbUploadTest)
	r.GET("/db/test/find", iodb.DbFindTest)
	r.POST("/db/test/find", iodb.DbUpdateTest)

	// Authentication
	r.POST("/login", authentication.Login)
	r.POST("/renew", authentication.RenewToken)
	r.POST("/logout", authentication.Logout)

	// return handler router to main()
	return r
}
