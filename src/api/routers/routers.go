package routers

import (
	"github.com/JosephS11723/CooPIR/src/api/handlers/debug"
	"github.com/JosephS11723/CooPIR/src/api/handlers/iodb"
	"github.com/JosephS11723/CooPIR/src/api/handlers/iojobserver"
	"github.com/JosephS11723/CooPIR/src/api/handlers/ioseaweed"
	"github.com/JosephS11723/CooPIR/src/api/lib/jobservertypes"
	"github.com/gin-gonic/gin"
)

func InitRouter() (*gin.Engine, chan jobservertypes.NewJob, chan jobservertypes.WorkerInfo) {
	// intialize engine with default middleware (TODO: replace)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// set low memory limit for multipart forms (8 MiB)
	r.MaxMultipartMemory = 8 << 20

	//create the channels that will be used for handlers to send info to the job queue
	job_channel := make(chan jobservertypes.NewJob)
	worker_channel := make(chan jobservertypes.WorkerInfo)

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

	//JOB QUEUE
	r.POST("/job/new", iojobserver.CreateNewJob(job_channel))

	// TOKEN AUTHENTICATION HANDLING

	// return handler router to main()
	return r, job_channel, worker_channel
}
