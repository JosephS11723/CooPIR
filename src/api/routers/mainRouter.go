package routers

import (
	"github.com/JosephS11723/CooPIR/src/api/handlers/agenthandler"
	"github.com/JosephS11723/CooPIR/src/api/handlers/authentication"
	"github.com/JosephS11723/CooPIR/src/api/handlers/debug"
	"github.com/JosephS11723/CooPIR/src/api/handlers/iodb"
	"github.com/JosephS11723/CooPIR/src/api/handlers/iojobs"
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
		AllowOrigins:     []string{"http://localhost:4200"},
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
	v1.GET("/file", ioseaweed.SWGETQuery)
	v1.GET("/file/:caseuuid/:fileuuid", ioseaweed.SWGETPath)
	v1.POST("/file", ioseaweed.SWPOST)
	v1.DELETE("/file", ioseaweed.SWDELETE)

	// MONGO-DB Tests
	// v1.GET("/test", iodb.DbPingTest)
	// v1.POST("/test", iodb.DbUploadTest)
	// v1.GET("/test/find", iodb.DbFindTest)
	// v1.POST("/test/find", iodb.DbUpdateTest)

	// MONGO-DB Final
	//v1.POST("/case/add", iodb.DbCreateCase)
	//v1.POST("/case/update", iodb.DbUpdateCase)
	//v1.GET("/case/find", iodb.DbGetCaseInfo)

	//MONGO-DB

	v1.GET("/case", iodb.DbGetCaseInfo)
	v1.POST("/case/new", iodb.DbCreateCase)
	v1.POST("/case/update", iodb.DbUpdateCase)
	v1.GET("/user", iodb.DbGetUserInfo)
	v1.POST("/user/new", iodb.DbCreateUser)
	v1.GET("/user/edit", iodb.GetUserEditUser)
	v1.POST("/user/update", iodb.DbUpdateUser)
	v1.GET("/cases", iodb.GetUserViewCases)
	v1.GET("/case/make", iodb.GetUserMakeCase)
	v1.GET("/file/info", iodb.GetFileInfo)
	v1.GET("/case/files", iodb.GetCaseFiles)

	// Authentication
	v1.POST("/auth/renew", authentication.RenewToken)
	v1.POST("/auth/logout", authentication.Logout)

	// group for job server requests
	v2 := r.Group("/api/v1/jobs")

	// get status of a job
	v2.GET("/status", iojobs.GetStatus)

	// submit a job
	v2.POST("/new", iojobs.CreateJob)

	// get work
	v2.GET("/getwork", iojobs.GetWork)

	// get file for work
	//v2.GET("/:jobuuid/:fileuuid", iojobs.GetWorkFile)

	// submit work
	v2.POST("/:jobuuid/result", iojobs.SubmitWork)

	v2.POST("/worker/new", iojobs.CreateWorker)

	// get job types based on available work
	v2.GET("/types", iojobs.GetAvailableJobTypes)

	// group for unauthenticated actions
	v3 := r.Group("/api/v1/auth")

	// login
	v3.POST("/login", authentication.Login)

	// log collection agent path
	v4 := r.Group("/api/v1/agent")
	v4.GET("/ws", agenthandler.AgentHandler)
	v4.GET("/all", agenthandler.GetAgents)
	v4.POST("/task", agenthandler.SubmitWork)

	// return handler router to main()
	return r
}
