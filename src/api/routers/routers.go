package routers

import (
	"github.com/JosephS11723/CooPIR/src/api/handlers/debug"
	"github.com/JosephS11723/CooPIR/src/api/handlers/iodb"
	"github.com/JosephS11723/CooPIR/src/api/handlers/iojobserver"
	"github.com/JosephS11723/CooPIR/src/api/handlers/ioseaweed"
	"github.com/gin-gonic/gin"
)

func InitRouter() (*gin.Engine, chan interface{}) {
	// intialize engine with default middleware (TODO: replace)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// set low memory limit for multipart forms (8 MiB)
	r.MaxMultipartMemory = 8 << 20

	//create the channel that will be used for handlers to send info to the job queue
	queue_chan := make(chan interface{})

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
	r.POST("/job/new", iojobserver.CreateNewJob(queue_chan))

	// TOKEN AUTHENTICATION HANDLING

	// return handler router to main()
	return r, queue_chan
}
