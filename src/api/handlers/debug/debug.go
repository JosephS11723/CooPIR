package debug

import (
	"log"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/api/lib/logtypes"
	"github.com/gin-gonic/gin"
)

// Response to ping-pong challenge
func PingPong(c *gin.Context) {
	// log ping pong request
	_, err := dbInterface.MakeCaseLog(c, "", c.MustGet("identity").(string), dbtypes.Info, logtypes.PingPong, nil)

	if err != nil {
		// failed to log
		log.Panicln("INTERNAL SERVER ERROR: LOG FILE CREATION FAILED")
	}

	c.JSON(200, gin.H{
		"data": "pong",
	})
}
