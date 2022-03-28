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

func DbGetUserInfo(c *gin.Context) dbtypes.Case {

	var json_request map[string]interface{}

	err := c.BindJSON(&json_request)

	if err != nil {
		log.Panicln(err)
	}

	var result = dbInterface.FindDocByFilter("Users", "UserMetadata", json_request)

	var dbCase dbtypes.Case

	err = result.Decode(&dbCase)

	if err != nil {
		log.Panicln(err)
	}

	return dbCase

}

func DbCreateUser(c *gin.Context) {

	var json_request dbtypes.NewUser

	err := c.BindJSON(&json_request)

	if err != nil {
		log.Panicln(err)
	}

	_, err = dbInterface.MakeUser(json_request)

	if err != nil {
		log.Panicln(err)
	}
}

func DbUpdateUser(c *gin.Context) {

	var json_request dbtypes.UpdateDoc

	c.BindJSON(&json_request)

	dbInterface.UpdateUser("Users", "UserMetadata", json_request)

}
