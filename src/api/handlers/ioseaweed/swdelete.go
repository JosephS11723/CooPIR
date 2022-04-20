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

// SWDELETE deletes a file from seaweedfs
func SWDELETE(c *gin.Context) {
	// get filename
	filename, success := c.GetQuery("fileuuid")

	// error if filename not provided
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no fileuuid provided"})
	}

	caseUUID, success := c.GetQuery("caseuuid")
	// error if caseuuid not provided
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no caseuuid provided"})
		return
	}

	// log file deletion
	_, err := dbInterface.MakeCaseLog(c, caseUUID, c.MustGet("identity").(string), dbtypes.Info, logtypes.FileDeleteAttempt, nil)
	if err != nil {
		// failed to log
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
	}

	// TODO: verify user is authorized to delete file  and case files are marked as editable (not likely)

	// log file deletion
	_, err = dbInterface.MakeCaseLog(c, caseUUID, c.MustGet("identity").(string), dbtypes.Info, logtypes.FileDelete, nil)
	if err != nil {
		// failed to log
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
	}

	// run delete function from lib
	err = swi.DELETEFile(c,filename, caseUUID)
	if err != nil {
		log.Panicln("INTERNAL SERVER ERROR: DELETE FAILED")
	}
}
