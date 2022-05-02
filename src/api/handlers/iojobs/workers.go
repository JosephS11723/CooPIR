package iojobs

import (
	"net/http"

	"github.com/JosephS11723/CooPIR/src/api/config"
	"github.com/JosephS11723/CooPIR/src/api/lib/coopirutil"
	"github.com/JosephS11723/CooPIR/src/api/lib/crypto"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/api/lib/logtypes"
	"github.com/gin-gonic/gin"
)

//used for creating a new job worker
func CreateWorker(c *gin.Context) {
	query_params := []string{"name", "jobtype"}

	single, _, err := coopirutil.ParseParams(query_params, c.Request.URL.Query())

	if err != nil {

		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

	}

	NewWorker := dbtypes.NewWorker{
		Name:    single["name"],
		JobType: single["jobtype"],
	}

	uuid, err := dbInterface.MakeWorker(NewWorker)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
	}

	// create token for worker
	token, err := crypto.CreateToken(uuid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	// worker creation log
	dbInterface.MakeCaseLog(c, logtypes.LogDefaultCaseUUID, uuid, dbtypes.Info, logtypes.WorkerStart, nil)

	// set token in cookie
	c.SetCookie("token", token, 9999999, "", "", false, config.HTTPOnly)

	c.JSON(200, gin.H{"uuid": uuid})
}

func CancelWorker() {

}

func GetWorkerInfo() {

}

func ReceiveResult() {

}
