package routers

import (
	"github.com/JosephS11723/CooPIR/src/api/handlers/debug"
	"github.com/JosephS11723/CooPIR/src/api/handlers/ioftp"
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

	// FTP TEST CODE
	r.GET("/file/:filename", ioftp.FtpDownload)
	r.POST("/file/:filename", ioftp.FtpUpload)
	r.DELETE("/file/:filename", ioftp.FtpDelete)

	// MONGO-DB

	// TOKEN AUTHENTICATION HANDLING

	// return handler router to main()
	return r
}
