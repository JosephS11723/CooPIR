package iodb

import (

	//"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/gin-gonic/gin"
	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
)

func DbGetCaseInfo(c *gin.Context) {

	var json_request map[string]interface{}

	c.BindJSON(&json_request)

	dbInterface.FindDocByFilter("Cases", "CaseMetadata", json_request)

}

func DbCreateCase(c *gin.Context) {

	var json_request dbtypes.Case

	c.BindJSON(&json_request)

	dbInterface.MakeCase(json_request)

}

func DbUpdateCase(c *gin.Context) {

	var json_request dbtypes.UpdateCase

	c.BindJSON(&json_request)

	dbInterface.UpdateCase("Cases", "CaseMetadata", json_request)

	//dbInterface.MakeCase(json_request)

}
