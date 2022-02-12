package iodb

import (
	"log"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DbPingTest(c *gin.Context) {
	log.Println("[DEBUG] iodb.DbTest()")
}

func DbUploadTest(c *gin.Context) {
	log.Println("[DEBUG] iodb.DbUploadTest()")

	var dbName string = "Users"
	var dbCollection string = "User"
	var result *mongo.InsertOneResult

	// Create a new user
	var testUser dbtypes.User = dbInterface.MakeUser("1", "testuser", "test@test.com", "supervisor", []string{"testcase", "thiscasedoesnotexist"}, "password")

	// for loop to add 3 of the same users
	for i := 0; i < 3; i++ {
		result = dbInterface.DbSingleInsert(dbName, dbCollection, testUser)

		log.Printf("[DEBUG] Inserted user document with _id: %v\n", result.InsertedID)
	}

	// Set to Cases db
	dbName = "Cases"

	// Create a new case
	dbCollection = "Case"
	var testCase dbtypes.Case = dbInterface.MakeCase("1", "testcase", "1/1/1976", "responder", "supervisor", []string{"testuser, anotheruser, Brandon"})

	result = dbInterface.DbSingleInsert(dbName, dbCollection, testCase)

	log.Printf("[DEBUG] Inserted case document with _id: %v\n", result.InsertedID)

	// Create a new file
	dbCollection = "File"
	var testfile dbtypes.File = dbInterface.MakeFile("1", "d5bb3ed1ccde75691e54f8f2da83a2fbf7eb9f0891ea141e67dd7f2b889ac479", "testfile", "case1", "test/dir", "1/1/1976", "supervisor", "admin")

	result = dbInterface.DbSingleInsert(dbName, dbCollection, testfile)

	log.Printf("[DEBUG] Inserted file document with _id: %v\n", result.InsertedID)

	// Create a new log
	dbCollection = "Log"
	var testLog dbtypes.Access = dbInterface.MakeAccess("1", "testfile", "testuse", "1/1/1976")

	result = dbInterface.DbSingleInsert(dbName, dbCollection, testLog)

	log.Printf("[DEBUG] Inserted access document with _id: %v\n", result.InsertedID)

}

func DbFindTest(c *gin.Context) {
	log.Println("[DEBUG] iodb.DbFindTest()")

	var dbName string = "Users"
	var dbCollection string = "User"

	filter := bson.M{"email": "test@test.com"}

	// Find user by filter
	result := dbInterface.FindDocsByFilter(dbName, dbCollection, filter)

	for _, user := range result {
		// un marshal the document into a user
		// var user dbtypes.User
		// err := bson.Unmarshal(document.(byte), &user)

		log.Print("[DEBUG] Found user document with _id: ", user["name"], " ", user["email"], " ", user["role"], " ", user["cases"], " ", user["saltedhash"], "\n")
	}

}
