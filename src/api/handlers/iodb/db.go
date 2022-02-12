package iodb

import (
	"log"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

	var dbName string = "Users"
	var dbCollection string = "User"

	// Create a new user
	var testUser dbtypes.User = dbInterface.MakeUser("testuser", "test@test.com", "supervisor", []string{"testcase", "thiscasedoesnotexist"}, "password")

	result := dbInterface.DbSingleInsert(client, ctx, dbName, dbCollection, testUser)

	log.Printf("[DEBUG] Inserted user document with _id: %v\n", result.InsertedID)

	// Set to Cases db
	dbName = "Cases"

	// Create a new case
	dbCollection = "Case"
	var testCase dbtypes.Case = dbInterface.MakeCase("testcase", "1/1/1976", "responder", "supervisor", []string{"testuser, anotheruser, Brandon"})

	result = dbInterface.DbSingleInsert(client, ctx, dbName, dbCollection, testCase)

	log.Printf("[DEBUG] Inserted case document with _id: %v\n", result.InsertedID)

	// Create a new file
	dbCollection = "File"
	var testfile dbtypes.File = dbInterface.MakeFile("d5bb3ed1ccde75691e54f8f2da83a2fbf7eb9f0891ea141e67dd7f2b889ac479", "testfile", "case1", "test/dir", "1/1/1976", "supervisor", "admin")

	result = dbInterface.DbSingleInsert(client, ctx, dbName, dbCollection, testfile)

	log.Printf("[DEBUG] Inserted file document with _id: %v\n", result.InsertedID)

	// Create a new log
	dbCollection = "Log"
	var testLog dbtypes.Access = dbInterface.MakeAccess("testfile", "testuse", "1/1/1976")

	result = dbInterface.DbSingleInsert(client, ctx, dbName, dbCollection, testLog)

	log.Printf("[DEBUG] Inserted access document with _id: %v\n", result.InsertedID)

}

func DbFindTest(c *gin.Context) {
	log.Println("[DEBUG] iodb.DbFindTest()")

	client, ctx, cancel, err := dbInterface.DbConnect()
	if err != nil {
		log.Panicln(err)
	}

	defer dbInterface.DbClose(client, ctx, cancel)

	var dbName string = "Users"
	var dbCollection string = "User"

	filter := bson.M{"name": "testuser", "email": "test@test.com"}

	// Find user by filter
	result := dbInterface.FindDocsByFilter(client, ctx, dbName, dbCollection, filter)

	log.Printf("[DEBUG] Found %v \n", result)

}
