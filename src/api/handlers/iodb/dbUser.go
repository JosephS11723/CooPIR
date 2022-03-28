package iodb

import (

	//"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/gin-gonic/gin"
	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
)

func DbGetUserInfo(c *gin.Context) {

	var json_request map[string]interface{}

	c.BindJSON(&json_request)

	dbInterface.FindDocByFilter("Users", , json_request)

}

func DbCreateUser(c *gin.Context) {

		var json_request map[string]{}

		err := c.BindJSON(&json_request)



		dbInterface.MakeUser(json_request)

}

func DbUpdateUser(c *gin.Context) {

	/*
		var json_request dbtypes.UpdateCase

		c.BindJSON(&json_request)

		dbInterface.UpdateCase("Cases", "CaseMetadata", json_request)
	*/
}
