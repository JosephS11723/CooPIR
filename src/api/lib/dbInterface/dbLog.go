package dbInterface

import (
	"time"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"go.mongodb.org/mongo-driver/mongo"
)

func MakeCaseLog(uuid string, level dbtypes.ErrorLevel, logtype string, content interface{}) (*mongo.InsertOneResult, error) {

	var unix_time = time.Now().Unix()

	var NewLog = dbtypes.Log{
		Uuid:    uuid,
		Level:   level,
		Type:    logtype,
		Time:    unix_time,
		Content: content,
	}

	result, err := DbSingleInsert("Logs", "Logs", NewLog)

	return result, err
}
