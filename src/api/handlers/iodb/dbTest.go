package iodb

import (
	"log"
	"strconv"
	"time"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
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
	var err error

	var testUser dbtypes.NewUser

	// Create three testusers in the database
	for i := 0; i < 3; i++ {

		testUser = dbtypes.NewUser{
			Name:     "testuser",
			Email:    "test" + strconv.Itoa(i) + "@test.com",
			Role:     "supervisor",
			Cases:    []string{"testcase", "thiscasedoesnotexist"},
			Password: "password",
		}

		result, err := dbInterface.MakeUser(testUser) //"testuser", "test"+strconv.Itoa(i)+"@test.com", "supervisor", []string{"testcase", "thiscasedoesnotexist"}, "password")

		if err != nil {
			log.Panicln("[ERROR] Failed to create testuser: " + err.Error())
		}

		log.Printf("[DEBUG] Inserted user document with _id: %v\n", result.InsertedID)
	}

	// Create new case in db
	//result = dbInterface.MakeCase("testcase", "1/1/1976", "responder", "supervisor", []string{"testuser, anotheruser, Brandon"})

	var NewCase = dbtypes.Case{

		Name:          "testcase",
		Date_created:  "1/1/1976",
		View_access:   "responder",
		Edit_access:   "supervisor",
		Collaborators: []string{"testuser, anotheruser, Brandon"},
	}

	result, caseuuid, err := dbInterface.MakeCase(NewCase)

	if err != nil {
		log.Panicln("[ERROR] Failed to create testcase: " + err.Error())
	}

	log.Printf("[DEBUG] Inserted case document with _id: %v\n", result.InsertedID)

	// Create a new file in db
	result, err = dbInterface.MakeFile(
		"7840eee0-e760-424a-95aa-3cf42c628d3f",
		[]string{"81051bcc2cf1bedf378224b0a93e2877",
			"ba8ab5a0280b953aa97435ff8946cbcbb2755a27",
			"7eb70257593da06f682a3ddda54a9d260d4fc514f645237f5ca74b08f8da61a6",
			"1b302a2f1e624a5fb5ad94ddc4e5f8bfd74d26fa37512d0e5face303d8c40eee0d0ffa3649f5da43f439914d128166cb6c4774a7caa3b174d7535451eb697b5d"},
		[]string{"tag", "tag2"},
		"testfile",
		"case1",
		"test/dir",
		strconv.Itoa(int(time.Now().UnixMilli())),
		"supervisor",
		"admin")

	if err != nil {
		log.Panicln("[ERROR] Failed to create testfile: " + err.Error())
	}

	log.Printf("[DEBUG] Inserted file document with _id: %v\n", result.InsertedID)

	// Create a new log in db
	result, err = dbInterface.MakeAccess("7840eee0-e760-424a-95aa-3cf42c628d3f", caseuuid, "testuse", "1/1/1976")

	if err != nil {
		log.Panicln("[ERROR] Failed to create testlog: " + err.Error())
	}

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
	result, _ := dbInterface.FindDocsByFilter(dbName, dbCollection, filter)

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
