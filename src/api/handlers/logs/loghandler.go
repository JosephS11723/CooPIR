package logs

import (
	"io"
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

	// create reader and writer to send to client
	reader, writer := io.Pipe()

	// defer close reader
	defer reader.Close()

	// get logs from db
	go dbInterface.GetCaseLogs(c, caseUUID, writer)

	// send datafromreader
	c.DataFromReader(http.StatusOK, 99999999, "application/json", reader, nil)

	// status 200
	c.Status(http.StatusOK)
	return
}
