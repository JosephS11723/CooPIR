package dbInterface

import (
	"net/http"
	"time"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/JosephS11723/CooPIR/src/api/lib/logtypes"
)

func MakeCaseLog(c *gin.Context, uuid string, user string, level dbtypes.ErrorLevel, logtype string, content interface{}) (*mongo.InsertOneResult, error) {
	// if caseuuid is not provided, set it to the default
	if uuid == "" {
		uuid = logtypes.LogDefaultCaseUUID
	}

	// if user is not provided, set it to the default
	if user == "" {
		user = logtypes.LogDefaultUserUUID
	}

	// get the current time in milliseconds since epoch
	var unix_time = time.Now().UnixMilli()

	// create the log entry
	var NewLog = dbtypes.Log{
		Uuid:    uuid,
		User:    user,
		Level:   level,
		Type:    logtype,
		Time:    unix_time,
		Content: content,
	}

	// insert the log entry into the database
	result, err := DbSingleInsert("Logs", "Logs", NewLog)

	if err != nil {
		// failed to log
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "log failure"})
	}

	return result, err
}
