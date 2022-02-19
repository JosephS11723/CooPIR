package iodb

import (
	"log"
	"strconv"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// DbPingTest is a test function to ping the database.
func DbPingTest(c *gin.Context) {
	log.Println("[DEBUG] iodb.DbTest()")
}

// DbUploadTest is a test function to upload a document to the database.
func DbUploadTest(c *gin.Context) {
	log.Println("[DEBUG] iodb.DbUploadTest()")

	var result *mongo.InsertOneResult

	// Create three testusers in the database
	for i := 0; i < 3; i++ {
		result, err := dbInterface.MakeUser("testuser", "test"+strconv.Itoa(i)+"@test.com", "supervisor", []string{"testcase", "thiscasedoesnotexist"}, "password")
		if err != nil {
			log.Panicln("[ERROR] Failed to create testuser: " + err.Error())
		}

		log.Printf("[DEBUG] Inserted user document with _id: %v\n", result.InsertedID)
	}

	// Create new case in db
	result = dbInterface.MakeCase("testcase", "1/1/1976", "responder", "supervisor", []string{"testuser, anotheruser, Brandon"})

	log.Printf("[DEBUG] Inserted case document with _id: %v\n", result.InsertedID)

	// Create a new file in db
	result = dbInterface.MakeFile([]string{"d5bb3ed1ccde75691e54f8f2da83a2fbf7eb9f0891ea141e67dd7f2b889ac479"}, []string{"tag", "tag2"}, "testfile", "case1", "test/dir", "1/1/1976", "supervisor", "admin")

	log.Printf("[DEBUG] Inserted file document with _id: %v\n", result.InsertedID)

	// Create a new log in db
	result = dbInterface.MakeAccess("testfile", "testuse", "1/1/1976")

	log.Printf("[DEBUG] Inserted access document with _id: %v\n", result.InsertedID)

}

// DbFindTest is a test function to find a document in the database.
func DbFindTest(c *gin.Context) {
	log.Println("[DEBUG] iodb.DbFindTest()")

	// variables for test
	var dbName string = "Users"
	var dbCollection string = "User"

	// create filter
	filter := bson.M{"email": "test@test.com"}

	// find user by filter
	result := dbInterface.FindDocsByFilter(dbName, dbCollection, filter)

	// iterate through results
	for _, user := range result {
		// un marshal the document into a user
		// var user dbtypes.User
		// err := bson.Unmarshal(document.(byte), &user)

		log.Print("[DEBUG] Found user document with _id: ", user["name"], " ", user["email"], " ", user["role"], " ", user["cases"], " ", user["saltedhash"], "\n")
	}
}

func DbUpdateTest(c *gin.Context) {
	log.Println("[DEBUG] iodb.DbUpdateTest()")

	var dbName string = "Users"
	var dbCollection string = "User"
	var name string = "testuser2"

	filter := bson.M{"email": "test0@test.com"}
	update := bson.D{{"$set",
		bson.D{
			{"name", name},
		},
	}}

	result := dbInterface.UpdateDoc(dbName, dbCollection, filter, update)

	log.Printf("[DEBUG] Updated user document with _id: %v\n", result.ModifiedCount)
}
