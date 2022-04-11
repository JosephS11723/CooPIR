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

// Requires json to have caseUUID field in request body
// Returns all case metadata
func DbGetCaseInfo(c *gin.Context) {
	var json_request map[string]interface{}

	err := c.BindJSON(&json_request)

	if err != nil {
		log.Panicln(err)
	}

	var caseUUID = json_request["uuid"].(string)

	//dbInterface.FindCase("Case", "CaseMetadata", json_request)
	result, err := dbInterface.FindDocByFilter("Cases", "CaseMetadata", bson.M{"uuid": caseUUID})

	if err != nil {
		c.AbortWithError(http.StatusNotFound, errors.New("file not found"))
	}

	var dbCase dbtypes.Case

	err = result.Decode(&dbCase)

	if err != nil {
		log.Panicln(err)
	}

	c.JSON(http.StatusOK, gin.H{"case": dbCase})
}

func DbCreateCase(c *gin.Context) {

	var json_request dbtypes.Case

	err := c.BindJSON(&json_request)

	if err != nil {
		log.Panicln(err)
	}

	// call make case (it does the sanity checks for us)
	_, caseUUID, err := dbInterface.MakeCase(json_request)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("case already exists"))
	}

	// send ok
	c.String(http.StatusOK, caseUUID)
}

func DbUpdateCase(c *gin.Context) {

	var json_request dbtypes.UpdateDoc

	err := c.BindJSON(&json_request)

	if err != nil {
		log.Panicln(err)
	}

	// TODO: check to see if case name already taken

	dbInterface.UpdateCase("Cases", "CaseMetadata", json_request)

	// send ok
	c.JSON(http.StatusOK, gin.H{"message": "Case updated"})
}

func GetUserViewCases(c *gin.Context) {

	var uuid = c.GetString("identity")

	var cases = dbInterface.UserCases(uuid)

	c.JSON(http.StatusOK, gin.H{"cases": cases})
}
