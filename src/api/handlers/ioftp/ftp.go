package ioftp

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/JosephS11723/CooPIR/src/api/lib/ftpInterface"
	"github.com/gin-gonic/gin"
)

func FtpUpload(c *gin.Context) {
	// connect to server
	client := ftpinterface.FtpConnect()

	// ensure connection close
	defer ftpinterface.FtpClose(client)

	// get file multipart stream
	file, err := c.FormFile("Filename")
	if err != nil {
		log.Fatal(err)
	}

	var filestream multipart.File

	filestream, err = file.Open()
	if err != nil {
		log.Fatal(err)
	}

	// get hash for filename
	// THIS IS DEBUG CODE FOR NOW. JUST LET IT BE
	// DO NOT LET THIS STAY IN THE FINAL RELEASE. IT IS NOT PRACTICAL AND WILL GET IN THE WAY
	filename := "test.txt"

	// upload file
	err = ftpinterface.WriteFile(client, filename, filestream)
	if err != nil {
		// upload failed
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("'%s' upload failed", file.Filename))
	} else {
		// upload succeeded
		c.String(http.StatusOK, fmt.Sprintf("'%s' upload succeeded", file.Filename))
	}
}

func FtpDownload(c *gin.Context) {
	// connect to server
	client := ftpinterface.FtpConnect()

	// ensure connection close
	defer ftpinterface.FtpClose(client)
}

func FtpDelete(c *gin.Context) {
	// connect to server
	client := ftpinterface.FtpConnect()

	// ensure connection close
	defer ftpinterface.FtpClose(client)
}
