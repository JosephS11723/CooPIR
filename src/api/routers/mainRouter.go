package routers

import (
	"github.com/JosephS11723/CooPIR/src/api/handlers/authentication"
	"github.com/JosephS11723/CooPIR/src/api/handlers/debug"
	"github.com/JosephS11723/CooPIR/src/api/handlers/iodb"
	"github.com/JosephS11723/CooPIR/src/api/handlers/ioseaweed"
	authmw "github.com/JosephS11723/CooPIR/src/api/middleware/authentication"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitMainRouter() *gin.Engine {
	// intialize engine with default middleware (TODO: replace)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS for cross-site whitelisting
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "access-control-allow-origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// set low memory limit for multipart forms (8 MiB)
	r.MaxMultipartMemory = 8 << 20

	// setup base path for api version v1
	v1 := r.Group("/api/v1")

	// authentication middleware
	v1.Use(authmw.AuthenticationMiddleware)

	// DEBUG REQUESTS
	// debug ping challenge
	v1.GET("/ping", debug.PingPong)

	// seaweedfs file storage transfer routes
	v1.GET("/file", ioseaweed.SWGET)
	v1.POST("/file", ioseaweed.SWPOST)
	v1.DELETE("/file", ioseaweed.SWDELETE)

	// MONGO-DB Tests
	v1.GET("/db/test", iodb.DbPingTest)
	v1.POST("/db/test", iodb.DbUploadTest)
	v1.GET("/db/test/find", iodb.DbFindTest)
	v1.POST("/db/test/find", iodb.DbUpdateTest)

	// MONGO-DB Final
	v1.POST("/db/case/add", iodb.DbCreateCase)
	v1.POST("/db/case/update", iodb.DbUpdateCase)
	v1.GET("db/case/find", iodb.DbGetCaseInfo)

	//MONGO-DB
	/*
		v1.GET("/db/case", iodb.DbGetCaseInfo)
		v1.POST("/db/case/new", iodb.DbCreateCase)
		v1.POST("/db/case/update", iodb.DbUpdateCase)
		v1.GET("/db/user", iodb.DbGetUserInfo)
		v1.POST("/db/user", iodb.Db)
	*/

	// Authentication
	v1.POST("/renew", authentication.RenewToken)
	v1.POST("/logout", authentication.Logout)

	// group for job server requests
	/*v2 := r.Group("/api/v1/jobs")

	// get status of a job
	v2.GET("/status", jobs.GetInfo)

	// submit a job
	v2.POST("new", jobs.CreateJob)

	// get work
	v2.GET("/work", jobs.GetWork)

	// submit work
	v2.POST("/work", jobs.SubmitWork)

	// get job results
	v2.GET("/results", jobs.GetResults)*/

	// group for logging in
	v3 := r.Group("/api/v1/auth")
	
	// login
	v3.POST("/login", authentication.Login)

	// return handler router to main()
	return r
}
