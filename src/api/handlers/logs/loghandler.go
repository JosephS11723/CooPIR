package logs

import (
	"log"
	"net/http"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/gin-gonic/gin"
)

func GetCaseLogs(c *gin.Context) {
	// get caseuuid from query
	caseUUID := c.Query("caseuuid")

	// check if empty
	if caseUUID == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "no case uuid in the query",
			},
		)
		log.Println("no case uuid in the query")
		return
	}

	// get logs from db
	logs, err := dbInterface.GetCaseLogs(caseUUID)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		log.Println(err.Error())
		return
	}

	// return logs
	c.JSON(
		http.StatusOK,
		gin.H{
			"logs": logs,
		},
	)
}
