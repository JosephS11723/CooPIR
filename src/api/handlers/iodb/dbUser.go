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

func DbGetUserInfo(c *gin.Context) dbtypes.User {

	var json_request map[string]interface{}

	err := c.BindJSON(&json_request)

	if err != nil {
		log.Panicln(err)
	}

	var result = dbInterface.FindDocByFilter("Users", "UserMetadata", json_request)

	var dbUser dbtypes.User

	err = result.Decode(&dbUser)

	if err != nil {
		log.Panicln(err)
	}

	dbUser.SaltedHash = "no password hash for you ;)"

	return c.JSON(http.StatusOK, gin.H{"user": dbUser})

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
