package ioseaweed

import (
	"log"
	//"mime/multipart"
	"net/http"
	swi "github.com/JosephS11723/CooPIR/src/api/lib/seaweedInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/security"
	"github.com/gin-gonic/gin"
	mongoUUID "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

func SWUpload(c *gin.Context) {
	// verify token
	if (!security.VerifyToken(c)) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "invalid token"})
		log.Panicln("INVALID TOKEN")
	}

	// get file multipart stream
	filestream, _, err := c.Request.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no file received"})
		log.Println(err)
		return
	}
	
	// set filename to randomly generated name. change after hash operation
	filename, err := mongoUUID.New()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to allocate temporary filename"})
		log.Println(err)
		return
	}

	// get hash for filename
	// THIS IS DEBUG CODE FOR NOW. JUST LET IT BE
	// DO NOT LET THIS STAY IN THE FINAL RELEASE. IT IS NOT PRACTICAL AND WILL GET IN THE WAY
	//filename := "test.txt"

	// upload file
	//err = ftpinterface.WriteFile(client, filename, filestream)
	err = swi.UploadFile(string(filename), filestream)
	if err != nil {
		// upload failed
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "upload error"})
		log.Panicln(err)
	} else {
		// upload succeeded
		c.String(http.StatusOK, "upload success")
		log.Panicln(err)
	}
}

func SWDownload(c *gin.Context) {
	// verify token
	if (!security.VerifyToken(c)) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "invalid token"})
		log.Panicln("INVALID TOKEN")
	}

	// TODO: get filename

	// download file
	err := swi.GetFile(filename, c)

	if err != nil {
		
	}

	response, err := http.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png")
		

		reader := response.Body
 		defer reader.Close()
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")

		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="gopher.png"`,
		}

		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}

func SWDelete(c *gin.Context) {
	// verify token
	if (!security.VerifyToken(c)) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "invalid token"})
		log.Panicln("INVALID TOKEN")
	}
	// connect to server
	client := ftpinterface.FtpConnect()

	ftpinterface.DeleteFile(client, "test.txt")

	// ensure connection close
	defer ftpinterface.FtpClose(client)
}
