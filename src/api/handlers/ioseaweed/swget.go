package ioseaweed

import (
	"log"
	"net/http"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/api/lib/logtypes"
	swi "github.com/JosephS11723/CooPIR/src/api/lib/seaweedInterface"
	"github.com/gin-gonic/gin"
)

// SWGETQuery gets a file from the seaweedfs server and returns it to the client with query parameters
func SWGETQuery(c *gin.Context) {
	// get filename
	filename, success := c.GetQuery("filename")

	// error if filename not provided
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no filename provided"})
		return
	}

	caseUUID, success := c.GetQuery("caseuuid")
	// error if caseuuid not provided
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no caseuuid provided"})
		return
	}

	// get file for user
	SWGET(c, filename, caseUUID)
}

// SWGETPath gets a file from the seaweedfs server and returns it to the client using the path parameters
func SWGETPath(c *gin.Context) {
	// get filename
	filename, success := c.Params.Get("filename")

	// error if filename not provided
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no filename provided"})
		return
	}

	caseUUID, success := c.Params.Get("caseuuid")
	// error if caseuuid not provided
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no caseuuid provided"})
		return
	}

	// get file for user
	SWGET(c, filename, caseUUID)
}

// SWGET gets a file from the seaweedfs server and returns it to the client
func SWGET(c *gin.Context, filename string, caseUUID string) {
	// log file download
	_, err := dbInterface.MakeCaseLog(c, caseUUID, c.MustGet("identity").(string), dbtypes.Info, logtypes.FileDownloadAttempt, filename)
	if err != nil {
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
	}

	// TODO: verify user is authorized to download file

	// download file through lib
	err = swi.GETFile(filename, caseUUID, c)

	// internal server error: failed to retrieve file data
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve data"})
		return
	}

	// log file download
	_, err = dbInterface.MakeCaseLog(c, caseUUID, c.MustGet("identity").(string), dbtypes.Info, logtypes.FileDownload, filename)
	if err != nil {
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
	}

	c.Status(http.StatusOK)
}
