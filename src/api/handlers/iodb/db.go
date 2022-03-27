package iodb

import (

	//"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"log"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
)

func DbGetCaseInfo(c *gin.Context) {

}

func DbCreateCase(c *gin.Context) {

	var json_request dbtypes.Case

	c.BindJSON(&json_request)

	dbInterface.MakeCase(json_request)

}

func DbUpdateCase(c *gin.Context) {

	var json_request interface{}

	//this is NOT the correct code for
	byte_data, err := c.GetRawData()

	if err != nil {

		log.Panicln(err)

	}

	bson.UnmarshalExtJSON(byte_data, true, json_request)

	filter := bson.M{"email": "test0@test.com"}

	/*
		update := bson.D{{"$set",
			bson.D{
				{"name", name},
			},
		}}
	*/

	dbInterface.UpdateDoc("Cases", "CaseMetadata", filter, update)

	//dbInterface.MakeCase(json_request)

}
