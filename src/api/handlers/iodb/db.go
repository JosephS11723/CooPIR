package iodb

import (
	"log"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/gin-gonic/gin"
)

func DbPingTest(c *gin.Context) {
	log.Println("[DEBUG] iodb.DbTest()")
	// Get Client, Context, CalcelFunc and
	// err from connect method.
	client, ctx, cancel, err := dbInterface.DbConnect()
	if err != nil {
		log.Panicln(err)
	}

	// Release resource when the main
	// function is returned.
	defer dbInterface.DbClose(client, ctx, cancel)

	// Ping mongoDB with Ping method
	dbInterface.DbPing(client, ctx)
}

func DbUploadTest(c *gin.Context) {
	log.Println("[DEBUG] iodb.DbUploadTest()")

	client, ctx, cancel, err := dbInterface.DbConnect()
	if err != nil {
		log.Panicln(err)
	}

	defer dbInterface.DbClose(client, ctx, cancel)

	var dbName string = "testingdb"
	var dbCollection string = "Users"

	var testUser dbtypes.User = dbtypes.User{
		UUID:  "testuser",
		Name:  "testuser",
		Email: "example@example.com",
		Role:  "admin",
		Cases: []string{"testcase1, testcase2"},
		Auth: []dbtypes.Authentication{
			{
				Salt: "salt",
				Pass: "pass",
			},
		},
	}

	result := dbInterface.DbSingleInsert(client, ctx, dbName, dbCollection, testUser)

	log.Printf("[DEBUG] Inserted document with _id: ObjectID(%v)\n", result.InsertedID)

}
