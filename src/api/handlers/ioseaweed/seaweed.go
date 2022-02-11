package ioseaweed

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/JosephS11723/CooPIR/src/api/config"
	libcrypto "github.com/JosephS11723/CooPIR/src/api/lib/crypto"
	swi "github.com/JosephS11723/CooPIR/src/api/lib/seaweedInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/security"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SWGET gets a file from the seaweedfs server and returns it to the client
func SWGET(c *gin.Context) {
	// verify token
	if !security.VerifyToken(c) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "invalid token"})
		log.Panicln("INVALID TOKEN")
	}

	// get filename
	filename, success := c.GetQuery("filename")

	// error if filename not provided
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no filename provided"})
	}

	// TODO: verify user is authorized to download file

	// download file through lib
	err := swi.GETFile(filename, c)

	// internal server error: failed to retrieve file data
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve data"})
	}
	c.Status(http.StatusOK)
}

// SWPOST uploads a file to seaweedfs from the client multipart form
func SWPOST(c *gin.Context) {
	var err error

	// verify token
	if !security.VerifyToken(c) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "invalid token"})
		log.Panicln("INVALID TOKEN")
	}

	// get file multipart stream
	filestream, _, err := c.Request.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no file received"})
		log.Panicln(err)
	}

	// set filename to randomly generated name. change after hash operation
	filename := uuid.New().String()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to allocate temporary filename"})
		log.Panicln(err)
	}

	// TODO: check db to make sure uuid does not already exist

	// create pipes
	md5Reader, md5Writer := io.Pipe()
	sha1Reader, sha1Writer := io.Pipe()
	sha256Reader, sha256Writer := io.Pipe()
	sha512Reader, sha512Writer := io.Pipe()
	POSTReader, POSTWriter := io.Pipe()

	// make channel for errors
	errChan := make(chan error, config.UpTaskCount)
	defer close(errChan)

	// make channel for hash result
	hashmd5Chan := make(chan []byte, 1)
	defer close(hashmd5Chan)

	// make channel for hash result
	hashsha1Chan := make(chan []byte, 1)
	defer close(hashsha1Chan)

	// make channel for hash result
	hashsha256Chan := make(chan []byte, 1)
	defer close(hashsha256Chan)

	// make channel for hash result
	hashsha512Chan := make(chan []byte, 1)
	defer close(hashsha512Chan)

	// sync group for synchronizing tasks
	var doWait sync.WaitGroup

	// addthe number of tasks to do
	doWait.Add(config.UpTaskCount)

	// spawn all the task goroutines. These look identical to
	// the TeeReader example, but pulled out into separate
	// methods for clarity
	go libcrypto.MD5FromReaderAsync(md5Reader, &doWait, errChan, hashmd5Chan)
	go libcrypto.Sha1FromReaderAsync(sha1Reader, &doWait, errChan, hashsha1Chan)
	go libcrypto.Sha256FromReaderAsync(sha256Reader, &doWait, errChan, hashsha256Chan)
	go libcrypto.Sha512FromReaderAsync(sha512Reader, &doWait, errChan, hashsha512Chan)
	go swi.POSTFile(filename, POSTReader, c.Copy(), &doWait, errChan)

	go func() {
		// after completing the copy, we need to close
		// the PipeWriters to propagate the EOF to all
		// PipeReaders to avoid deadlock
		defer md5Writer.Close()
		defer sha1Writer.Close()
		defer sha256Writer.Close()
		defer sha512Writer.Close()
		defer POSTWriter.Close()

		// build the multiwriter for all the pipes
		mw := io.MultiWriter(md5Writer, sha1Writer, sha256Writer, sha512Writer, POSTWriter)

		// copy the data into the multiwriter
		io.Copy(mw, filestream)
	}()

	// read hash output from functions
	filemd5Hash := <-hashmd5Chan
	filesha1Hash := <-hashsha1Chan
	filesha256Hash := <-hashsha256Chan
	filesha512Hash := <-hashsha512Chan

	// if logging for upload is enabled, print hashes to log
	if config.DoUploadLogging {
		fmt.Println("MD5:   ", hex.EncodeToString(filemd5Hash))
		fmt.Println("SHA1:  ", hex.EncodeToString(filesha1Hash))
		fmt.Println("SHA256:", hex.EncodeToString(filesha256Hash))
		fmt.Println("SHA512:", hex.EncodeToString(filesha512Hash))
	}

	// wait until all tasks are done
	doWait.Wait()

	// check errors
	for i := 0; i < config.UpTaskCount; i++ {
		err = <-errChan
		if err != nil {
			// upload failed, panic!
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "upload error"})
			log.Panicln(err)
		}
	}

	// TODO: check mongo for file existence and remove if duplicate

	// TODO: after checking if file exists, add new file information to the database

	// upload succeeded
	c.String(http.StatusOK, filename)
}

// SWDELETE deletes a file from seaweedfs
func SWDELETE(c *gin.Context) {
	// verify token
	if !security.VerifyToken(c) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "invalid token"})
		log.Panicln("INVALID TOKEN")
	}

	// get filename
	filename, success := c.GetQuery("filename")

	// error if filename not provided
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no filename provided"})
	}

	// TODO: verify user is authorized to delete file  and case files are marked as editable (not likely)

	// run delete function from lib
	err := swi.DELETEFile(filename, c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to delete"})
		log.Panicln("INTERNAL SERVER ERROR: DELETE FAILED")
	}
}
