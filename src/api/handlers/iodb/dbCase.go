package iodb

import (

	//"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"

	"log"

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

	dbInterface.FindDocByFilter("Cases", "CaseMetadata", json_request)

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

	var json_request dbtypes.UpdateCase

	err := c.BindJSON(&json_request)

	if err != nil {
		log.Panicln(err)
	}

	dbInterface.UpdateCase("Cases", "CaseMetadata", json_request)

}
