package ioftp

import (
	"log"
	//"mime/multipart"
	"net/http"

	"github.com/JosephS11723/CooPIR/src/api/lib/ftpInterface"
	"github.com/gin-gonic/gin"
)

func FtpUpload(c *gin.Context) {
	// get file multipart stream
	filestream, _, err := c.Request.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no file received",})
		log.Println(err)
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "upload error",})
		log.Println(err)
		return
	}

	// get hash for filename
	// THIS IS DEBUG CODE FOR NOW. JUST LET IT BE
	// DO NOT LET THIS STAY IN THE FINAL RELEASE. IT IS NOT PRACTICAL AND WILL GET IN THE WAY
	filename := "test.txt"

	// connect to server
	client := ftpinterface.FtpConnect()

	// ensure connection close
	defer ftpinterface.FtpClose(client)

	// upload file
	err = ftpinterface.WriteFile(client, filename, filestream)
	if err != nil {
		// upload failed
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "upload error",})
		log.Println(err)
	} else {
		// upload succeeded
		c.String(http.StatusOK, "upload succeeded")
		log.Println(err)
	}
}

func FtpDownload(c *gin.Context) {
	// connect to server
	client := ftpinterface.FtpConnect()

	// TEST CODE. SHOULD ONLY RETURN THE ONE FILE
	ftpinterface.ReadFile(client, "test.txt", c.Writer)

	// ensure connection close
	defer ftpinterface.FtpClose(client)
}

func FtpDelete(c *gin.Context) {
	// connect to server
	client := ftpinterface.FtpConnect()

	ftpinterface.DeleteFile(client, "test.txt")

	// ensure connection close
	defer ftpinterface.FtpClose(client)
}
