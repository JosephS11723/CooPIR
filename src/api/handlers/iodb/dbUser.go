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
	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
)

// Requires json to have userUUID field in request body
// Returns all user metadata
func DbGetUserInfo(c *gin.Context) {
	var json_request map[string]interface{}

	err := c.BindJSON(&json_request)

	if err != nil {
		log.Panicln(err)
	}

	var userUUID = json_request["userUUID"].(string)

	result, err := dbInterface.FindDocByFilter("Users", "UserMetadata", bson.M{"uuid": userUUID})

	if err != nil {
		// 404 user not found
		c.AbortWithError(http.StatusNotFound, errors.New("not found"))
	}

	var dbUser dbtypes.User

	err = result.Decode(&dbUser)

	if err != nil {
		log.Panicln(err)
	}

	dbUser.SaltedHash = "no password hash for you ;)"

	c.JSON(http.StatusOK, gin.H{"user": dbUser})

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

// Returns true if the user can edit or make cases
func GetUserMakeCase(c *gin.Context) {

	var uuid = c.GetString("identity")

	var allow = dbInterface.UserSupervisorPermission(uuid)

	c.JSON(http.StatusOK, gin.H{"allow": allow})

}

// Returns true if the user can edit or make users
func GetUserEditUser(c *gin.Context) {
	var uuid = c.GetString("identity")

	var allow = dbInterface.UserAdminPermission(uuid)

	c.JSON(http.StatusOK, gin.H{"allow": allow})
}

// Returns true if the user can edit or make users
func GetUserAddFile(c *gin.Context) {
	var uuid = c.GetString("identity")

	var allow = dbInterface.UserResponderPermission(uuid)

	c.JSON(http.StatusOK, gin.H{"allow": allow})
}
