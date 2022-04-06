package dbInterface

import (
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
)

func MakeCaseLog(uuid string, level Level, logtype string, content interface{}) (*mongo.InsertOneResult, error) {

	var unix_time = time.Now().Unix()

	var NewLog = Log {
		Uuid: uuid,
		Level: level,
		Type: logtype,
		Time: unix_time,
		Content: content,
	}

	result, err = DbSingleInsert("Logs", "Logs", NewLog)

	return result, err
}
