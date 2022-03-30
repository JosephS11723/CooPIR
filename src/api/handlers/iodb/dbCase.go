package iodb

import (

	//"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"

	"log"
	"net/http"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/gin-gonic/gin"
	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
)

func DbGetCaseInfo(c *gin.Context) {

	var json_request map[string]interface{}

	err := c.BindJSON(&json_request)

	if err != nil {
		log.Panicln(err)
	}

	//dbInterface.FindCase("Case", "CaseMetadata", json_request)
	var result = dbInterface.FindDocByFilter("Cases", "CaseMetadata", json_request)

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

	dbInterface.MakeCase(json_request)

}

func DbUpdateCase(c *gin.Context) {

	var json_request dbtypes.UpdateDoc

	err := c.BindJSON(&json_request)

	if err != nil {
		log.Panicln(err)
	}

	dbInterface.UpdateCase("Cases", "CaseMetadata", json_request)

}
