package iodb

import (

	//"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"

	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	httputil "github.com/JosephS11723/CooPIR/src/api/lib/coopirutil"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/api/lib/logtypes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// Requires json to have caseUUID field in request body
// Returns all case metadata
func DbGetCaseInfo(c *gin.Context) {

	/*var json_request map[string]interface{}

	err := c.BindJSON(&json_request)

	if err != nil {
		log.Panicln(err)
	}

	var caseUUID = json_request["uuid"].(string)*/

	caseUUID := c.Query("uuid")

	if caseUUID == "" {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": "query did not contain job uuid",
			},
		)
	}

	// log case info request
	_, err := dbInterface.MakeCaseLog(c, caseUUID, c.MustGet("identity").(string), dbtypes.Info, logtypes.GetCaseInfo, nil)

	if err != nil {
		// failed to log
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
	}

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
	// TODO: check if user is able to create cases

	/*
		var json_request dbtypes.Case

		err := c.BindJSON(&json_request)

		if err != nil {
			log.Panicln(err)
		}
	*/

	query_params := []string{"name", "description", "viewAccess", "editAccess", "collabs"}

	singles, multi, err := httputil.ParseParams(query_params, c.Request.URL.Query())

	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	newCase := dbtypes.Case{
		Name:          singles["name"],
		Description:   singles["description"],
		Date_created:  strconv.Itoa(int(time.Now().UnixMilli())),
		ViewAccess:    singles["viewAccess"],
		EditAccess:    singles["editAccess"],
		Collaborators: multi["collabs"],
	}

	// print newCase
	log.Printf("%+v\n", newCase)

	// log
	_, err = dbInterface.MakeCaseLog(c, "", c.MustGet("identity").(string), dbtypes.Info, logtypes.CreateCaseAttempt, nil)

	if err != nil {
		// failed to log
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
		return
	}

	// call make case (it does the sanity checks for us). // TODO: figure out where to put CreateCaseFailure log
	_, caseUUID, err := dbInterface.MakeCase(newCase)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("case already exists"))
		return
	}

	// log
	_, err = dbInterface.MakeCaseLog(c, caseUUID, c.MustGet("identity").(string), dbtypes.Info, logtypes.CreateCase, nil)

	if err != nil {
		// failed to log
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
		return
	}

	// send ok
	c.String(http.StatusOK, caseUUID)
}

func DbUpdateCase(c *gin.Context) {
	//TODO: the log functions require the case uuid that is being specified in the request body.
	var json_request dbtypes.UpdateDoc

	err := c.BindJSON(&json_request)

	if err != nil {
		log.Panicln(err)
	}

	// log
	_, err = dbInterface.MakeCaseLog(c, "", c.MustGet("identity").(string), dbtypes.Info, logtypes.UpdateCaseAttempt, nil)

	if err != nil {
		// failed to log
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
	}

	// TODO: check to see if case name already taken
	// TODO: add UpdateCaseFailure log if user does not have permission to update case info

	dbInterface.UpdateCase("Cases", "CaseMetadata", json_request)

	// log
	_, err = dbInterface.MakeCaseLog(c, "", c.MustGet("identity").(string), dbtypes.Info, logtypes.UpdateCase, nil)

	if err != nil {
		// failed to log
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
	}

	// send ok
	c.JSON(http.StatusOK, gin.H{"message": "Case updated"})
}

func GetUserViewCases(c *gin.Context) {
	var uuid = c.GetString("identity")

	var cases = dbInterface.UserCases(uuid)

	// log
	_, err := dbInterface.MakeCaseLog(c, "", c.MustGet("identity").(string), dbtypes.Info, logtypes.GetCases, gin.H{"cases": cases})

	if err != nil {
		// failed to log
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
	}

	c.JSON(http.StatusOK, gin.H{"cases": cases})
}
