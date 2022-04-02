package ioseaweed

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/JosephS11723/CooPIR/src/api/config"
	libcrypto "github.com/JosephS11723/CooPIR/src/api/lib/crypto"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	swi "github.com/JosephS11723/CooPIR/src/api/lib/seaweedInterface"
	"github.com/gin-gonic/gin"
)

// SWGET gets a file from the seaweedfs server and returns it to the client
func SWGET(c *gin.Context) {
	// get filename
	filename, success := c.GetQuery("filename")

	// error if filename not provided
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no filename provided"})
		return
	}

	caseUUID, success := c.GetQuery("caseuuid")
	// error if caseuuid not provided
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no caseuuid provided"})
		return
	}

	// TODO: verify user is authorized to download file

	// download file through lib
	err := swi.GETFile(filename, caseUUID, c)

	// internal server error: failed to retrieve file data
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve data"})
		return
	}
	c.Status(http.StatusOK)
}

// SWPOST uploads a file to seaweedfs from the client multipart form
func SWPOST(c *gin.Context) {
	var err error

	caseUUID, success := c.GetQuery("caseuuid")
	// error if caseuuid not provided
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no caseuuid provided"})
		return
	}

	originalFilename, success := c.GetQuery("filename")
	// error if filedir not provided
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no filedir provided"})
		return
	}

	// get file multipart stream
	filestream, _, err := c.Request.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no file received"})
		log.Panicln(err)
	}

	// ensure case name is valid
	fmt.Println(caseUUID)
	_, err = dbInterface.FindCaseNameByUUID(caseUUID)

	if err != nil {
		log.Println(err)
		// case does not exist
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid case UUID"})
		return
	}

	// TODO: ensure user is authorized to upload file to case

	// set filename to randomly generated name. change after hash operation
	// Use MakeUuid from dbInterface to ensure unique filename
	filename, err := dbInterface.MakeUuid()

	// error if failed to generate uuid
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to generate uuid"})
		return
	}

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
	go swi.POSTFile(filename, caseUUID, POSTReader, c.Copy(), &doWait, errChan)

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

	// check mongo for file existence and remove if duplicate
	_, err = dbInterface.FindFileByHash(hex.EncodeToString(filesha512Hash), caseUUID)
	
	if err == nil {
		// file already exists, remove it
		err = swi.DELETEFile(filename, caseUUID, c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to delete file"})
			log.Panicln(err)
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "file already exists"})
		return
	}

	_, err = dbInterface.MakeFile(
		filename,
		[]string{
			hex.EncodeToString(filemd5Hash),
			hex.EncodeToString(filesha1Hash),
			hex.EncodeToString(filesha256Hash),
			hex.EncodeToString(filesha512Hash),
		},
		[]string{},
		originalFilename,
		caseUUID,
		"/files/"+caseUUID+"/"+filename,
		time.Now().Local().String(),
		"supervisor",
		"admin",
	)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to create file"})
		log.Panicln(err)
	}

	// upload succeeded
	c.String(http.StatusOK, filename)
}

// SWDELETE deletes a file from seaweedfs
func SWDELETE(c *gin.Context) {
	// get filename
	filename, success := c.GetQuery("filename")

	// error if filename not provided
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no filename provided"})
	}

	caseUUID, success := c.GetQuery("caseuuid")
	// error if caseuuid not provided
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no caseuuid provided"})
		return
	}

	// TODO: verify user is authorized to delete file  and case files are marked as editable (not likely)

	// run delete function from lib
	err := swi.DELETEFile(filename, caseUUID, c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to delete"})
		log.Panicln("INTERNAL SERVER ERROR: DELETE FAILED")
	}
}
