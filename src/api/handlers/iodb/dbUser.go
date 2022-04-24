package iodb

import (

	//"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"

	"errors"
	"log"
	"net/http"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/api/lib/httputil"
	"github.com/JosephS11723/CooPIR/src/api/lib/logtypes"
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

	// log
	_, err = dbInterface.MakeCaseLog(c, "", c.MustGet("identity").(string), dbtypes.Info, logtypes.GetUserInfoAttempt, gin.H{"userUUID": userUUID})

	if err != nil {
		// failed to log
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
	}

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

	// TODO: why can't we just not return this field??? we can make a custom json from the db query
	dbUser.SaltedHash = "no password hash for you ;)"

	// log
	_, err = dbInterface.MakeCaseLog(c, "", c.MustGet("identity").(string), dbtypes.Info, logtypes.GetUserInfo, gin.H{"userUUID": userUUID})

	if err != nil {
		// failed to log
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
	}

	c.JSON(http.StatusOK, gin.H{"user": dbUser})

}

func DbCreateUser(c *gin.Context) {

	query_params := []string{"name", "email", "role", "cases"}

	singles, multi, err := httputil.ParseParams(query_params, c.Request.URL.Query())

	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
	}

	newUser := dbtypes.NewUser{
		Name:  singles["name"],
		Email: singles["email"],
		Role:  singles["role"],
		Cases: multi["cases"],
	}

	// log
	_, err = dbInterface.MakeCaseLog(c, "", c.MustGet("identity").(string), dbtypes.Info, logtypes.CreateUserAttempt, nil)

	if err != nil {
		// failed to log
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
	}

	_, err = dbInterface.MakeUser(newUser)

	// log
	_, err = dbInterface.MakeCaseLog(c, "", c.MustGet("identity").(string), dbtypes.Info, logtypes.CreateUser, nil)

	if err != nil {
		// failed to log
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
	}

	if err != nil {
		log.Panicln(err)
	}
}

//currently not implemented by client, so don't worry about it
func DbUpdateUser(c *gin.Context) {

	/*
		query_params := []string{"name", "email", "role", "cases"}

		singles, multi, err := httputil.ParseParams(query_params, c.Request.URL.Query())

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{
					"error": err.Error(),
				},
			)
		}
	*/

	// TODO: add target user and new information to log entries
	// TODO: add log attempt
	var json_request dbtypes.UpdateDoc

	c.BindJSON(&json_request)

	dbInterface.UpdateUser("Users", "UserMetadata", json_request)

	// log
	_, err := dbInterface.MakeCaseLog(c, "", c.MustGet("identity").(string), dbtypes.Info, logtypes.UpdateUser, nil)
	if err != nil {
		// failed to log
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
	}
}

// Returns true if the user can edit or make cases
func GetUserMakeCase(c *gin.Context) {
	var uuid = c.GetString("identity")

	var allow = dbInterface.UserSupervisorPermission(uuid)

	// log
	_, err := dbInterface.MakeCaseLog(c, "", c.MustGet("identity").(string), dbtypes.Info, logtypes.GetUserMakeCase, nil)
	if err != nil {
		// failed to log
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
	}

	c.JSON(http.StatusOK, gin.H{"allow": allow})

}

// Returns true if the user can edit or make users
func GetUserEditUser(c *gin.Context) {
	var uuid = c.GetString("identity")

	// log
	_, err := dbInterface.MakeCaseLog(c, "", c.MustGet("identity").(string), dbtypes.Info, logtypes.GetUserEditUser, nil)
	if err != nil {
		// failed to log
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
	}

	var allow = dbInterface.UserAdminPermission(uuid)

	c.JSON(http.StatusOK, gin.H{"allow": allow})
}

// Returns true if the user can edit or make users
func GetUserAddFile(c *gin.Context) {
	// TODO: add logs to this function
	var uuid = c.GetString("identity")

	// log
	_, err := dbInterface.MakeCaseLog(c, "", c.MustGet("identity").(string), dbtypes.Info, logtypes.GetUserAddFile, nil)

	if err != nil {
		// failed to log
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
	}

	var allow = dbInterface.UserResponderPermission(uuid)

	c.JSON(http.StatusOK, gin.H{"allow": allow})
}
