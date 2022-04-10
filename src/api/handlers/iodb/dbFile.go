package iodb

import (

	//"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"

	"errors"
	"log"
	"net/http"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// Takes the caseUUID and fileUUID from json_request and returns the file metadata
func GetFileInfo(c *gin.Context) {
	var json_request map[string]interface{}

	err := c.BindJSON(&json_request)

	if err != nil {
		log.Panicln(err)
	}

	var caseUUID string = json_request["caseUUID"].(string)
	var fileUUID string = json_request["fileUUID"].(string)

	//dbInterface.FindCase("Case", "CaseMetadata", json_request)
	result, err := dbInterface.FindDocByFilter("Cases", caseUUID, bson.M{"uuid": fileUUID})

	if err != nil {
		c.AbortWithError(http.StatusNotFound, errors.New("file not found"))
	}

	var dbFile dbtypes.File

	err = result.Decode(&dbFile)

	if err != nil {
		log.Panicln(err)
	}

	c.JSON(http.StatusOK, gin.H{"file": dbFile})
}

// Takes the case UUID and returns all the files in the case
func GetCaseFiles(c *gin.Context) {

	//var json_request map[string]interface{}

	//err := c.BindJSON(&json_request)

	var caseUUID string = c.Params.ByName("uuid")

	if caseUUID == "" {
		c.AbortWithError(http.StatusNotFound, errors.New("files not found"))
	}

	/*
		if err != nil {
			log.Panicln(err)
		}*/

	files, err := dbInterface.FindFilesByCase(caseUUID)

	if err != nil {
		c.AbortWithError(http.StatusNotFound, errors.New("files not found"))
	}

	c.JSON(http.StatusOK, gin.H{"files": files})

}
