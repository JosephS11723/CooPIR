package dbInterface

import (
	// "io"
	"context"
	"log"

	// _ "sync"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Find the user's Role with an UUID
func findUserRoleByUUID(uuid string) (string, error) {
	var dbName string = "Users"
	var dbCollection string = "UserMetadata"
	result, err := FindDocByFilter(dbName, dbCollection, bson.M{"uuid": uuid})

	if err != nil {
		return "", err
	}

	var dbUser dbtypes.User
	err = result.Decode(&dbUser)

	if err != nil {
		return "", err
	}

	return dbUser.Role, nil
}

// Find the user's email with an UUID
func FindUserEmailByUUID(uuid string) (string, error) {
	var dbName string = "Users"
	var dbCollection string = "UserMetadata"
	result, err := FindDocByFilter(dbName, dbCollection, bson.M{"uuid": uuid})

	if err != nil {
		return "", err
	}

	var dbUser dbtypes.User
	err = result.Decode(&dbUser)

	if err != nil {
		return "", err
	}

	return dbUser.Email, nil
}

// Find the user's UUID with an email
func FindUserUUIDByEmail(email string) (string, error) {
	var dbName string = "Users"
	var dbCollection string = "UserMetadata"
	result, err := FindDocByFilter(dbName, dbCollection, bson.M{"email": email})

	if err != nil {
		return "", err
	}

	var dbUser dbtypes.User
	err = result.Decode(&dbUser)

	return dbUser.UUID, nil
}

// Finds the case name from CaseMetadata collection using the case UUID.
func FindCaseNameByUUID(uuid string) (string, error) {
	var dbName string = "Cases"
	var dbCollection string = "CaseMetadata"
	result, err := FindDocByFilter(dbName, dbCollection, bson.M{"uuid": uuid})

	if err != nil {
		return "", err
	}

	var dbCase dbtypes.Case
	err = result.Decode(&dbCase)

	return dbCase.Name, err
}

// Finds the uuid of a case by the case name
func FindCaseUUIDByName(name string) (string, error) {
	var dbName string = "Cases"
	var dbCollection string = "CaseMetadata"
	result, err := FindDocByFilter(dbName, dbCollection, bson.M{"name": name})

	if err != nil {
		return "", err
	}

	var dbCase dbtypes.Case
	err = result.Decode(&dbCase)

	return dbCase.UUID, err
}

// Returns a list of Cases that a role can view
// Returns a slice of Case UUIDs or an error when no cases are found
func retrieveCasesByViewRole(role string) ([]string, error) {
	var dbName string = "Cases"
	var dbCollection string = "CaseMetadata"
	result, err := FindDocsByFilter(dbName, dbCollection, bson.M{})

	log.Println(result)

	if err != nil {
		return nil, err
	}

	var caseList []string
	// var Access dbtypes.AccessLevel.ToInt(role)
	var Access dbtypes.AccessLevel

	// TODO: Might not work
	for _, doc := range result {
		log.Println("View access:", doc["viewaccess"])
		log.Println("Access level:", Access.ToInt(role))
		if Access.ToInt(role) >= Access.ToInt(doc["viewaccess"].(string)) {
			caseList = append(caseList, doc["uuid"].(string))
		}
	}

	return caseList, nil
}

// Takes a user UUID and returns a slice of Case UUIDS the user can view
func RetrieveViewCasesByUserUUID(uuid string) ([]string, error) {
	role, err := findUserRoleByUUID(uuid)

	log.Println("USER ROLE", role)

	if err != nil {
		return nil, err
	}

	caseList, err := retrieveCasesByViewRole(role)

	if err != nil {
		return nil, err
	}

	return caseList, nil
}

func FindFileByHash(hash string, dbCollection string) (string, error) {
	var dbName string = "Cases"
	result, err := FindDocByFilter(dbName, dbCollection, bson.M{"sha512": hash})

	if err != nil {
		return "", err
	}

	var dbFile dbtypes.File
	err = result.Decode(&dbFile)

	return dbFile.UUID, err
}

// FindDocsByFilter finds multiple documents in a collection by a filter and
// RETURN a slice of documents (bson.m)
func FindDocsByFilter(dbname string, collection string, filter bson.M) ([]bson.M, error) {
	// connect to db
	client, ctx, cancel, err := dbConnect()

	// defer closing db connection
	defer dbClose(client, ctx, cancel)

	if err != nil {
		return nil, err
	}

	// get collection
	coll := client.Database(dbname).Collection(collection)

	// run find function and get cursor
	cur, err := coll.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	// create slice to hold documents
	var docList []bson.M

	// iterate through the documents
	for cur.Next(context.Background()) {
		// decode cur into a interface
		var doc bson.M
		err := cur.Decode(&doc)

		if err != nil {
			return nil, err
		}
		// log.Println("[DEBUG] internal result: ", doc)

		// append doc to docList
		docList = append(docList, doc)
	}

	// return list of documents
	return docList, nil
}

// FindDocByFilter finds a single document in a collection by a filter and RETURN a document (*mongo.SingleResult)
func FindDocByFilter(dbname string, collection string, filter bson.M) (*mongo.SingleResult, error) {
	// connect to db
	client, ctx, cancel, err := dbConnect()

	// defer closing db connection
	defer dbClose(client, ctx, cancel)

	if err != nil {
		log.Panicln(err)
	}

	// get the collection
	coll := client.Database(dbname).Collection(collection)

	// find the document
	var result *mongo.SingleResult = coll.FindOne(ctx, filter)

	// return the result
	return result, result.Err()
}

// Checks if Case name already exist in the mongo database
func DoesCaseExist(name string) (bool, error) {
	var dbName string = "Cases"
	var dbCollection string = "CaseMetadata"

	// connect to db
	client, ctx, cancel, err := dbConnect()
	if err != nil {
		log.Panicln(err)
	}

	defer dbClose(client, ctx, cancel)

	// get the collection
	coll := client.Database(dbName).Collection(dbCollection)

	// find the document
	var result *mongo.SingleResult = coll.FindOne(ctx, bson.M{"name": name})

	// return true if document exists
	return result.Err() != mongo.ErrNoDocuments, nil
}

// Finds all collections in the database and returns a slice of strings
// findCollections(bdName string) []string
func findCollections(dbName string) []string {
	// connect to db
	client, ctx, cancel, err := dbConnect()
	if err != nil {
		log.Panicln(err)
	}

	defer dbClose(client, ctx, cancel)

	// Obtain the DB, by name. db will have the type
	// *mongo.Database
	db := client.Database(dbName)

	// use a filter to only select all collections
	result, err := db.ListCollectionNames(ctx, bson.D{{"options.capped", false}})

	if err != nil {
		log.Panicln(err)
	}

	for _, coll := range result {
		log.Println(coll)
	}

	return result
}

// Find user role by user UUID
// func findRolesbyUUID(uuid string) string {
func FindRolebyUUID(uuid string) string {

	result, err := FindDocByFilter("Users", "UserMetadata", bson.M{"uuid": uuid})

	if err != nil {
		log.Panicln(err)
	}

	var dbUser dbtypes.User

	err = result.Decode(&dbUser)

	if err != nil {
		log.Panicln(err)
	}

	return dbUser.Role
}

// Get user uuid by email
// func findUUIDbyEmail(email string) string {
func FindUUIDbyEmail(email string) string {

	result, err := FindDocByFilter("Users", "UserMetadata", bson.M{"email": email})

	if err != nil {
		log.Panicln(err)
	}

	var dbUser dbtypes.User

	err = result.Decode(&dbUser)

	if err != nil {
		log.Panicln(err)
	}

	return dbUser.Role
}

// Check if UUID exists in the collection. Returns true if the document exists.
func doesUuidExist(dbname string, collection string, uuid string) (bool, error) {
	// connect to db
	client, ctx, cancel, err := dbConnect()
	if err != nil {
		log.Panicln(err)
	}

	defer dbClose(client, ctx, cancel)

	// get the collection
	coll := client.Database(dbname).Collection(collection)

	// find the document
	var result *mongo.SingleResult = coll.FindOne(ctx, bson.M{"uuid": uuid})

	// return true if document exists
	return result.Err() != mongo.ErrNoDocuments, nil
}

// Function to check if user email exists in the database.
// Returns true if the document exists.
func doesEmailExist(email string) (bool, error) {
	var dbName string = "Users"
	var dbCollection string = "UserMetadata"

	// connect to db
	client, ctx, cancel, err := dbConnect()
	if err != nil {
		log.Panicln(err)
	}

	defer dbClose(client, ctx, cancel)

	// get the collection
	coll := client.Database(dbName).Collection(dbCollection)

	// find the document
	var result *mongo.SingleResult = coll.FindOne(ctx, bson.M{"email": email})

	// return true if document exists
	return result.Err() != mongo.ErrNoDocuments, nil
}

func RetrieveHashByEmail(email string) (string, error) {
	var dbName string = "Users"
	var dbCollection string = "UserMetadata"
	var filter bson.M = bson.M{"email": email}

	result, err := FindDocByFilter(dbName, dbCollection, filter)

	if err != nil {
		return "", err
	}

	var userStructTmp dbtypes.User = dbtypes.User{}
	var hash string
	result.Decode(&userStructTmp)

	hash = userStructTmp.SaltedHash

	return hash, nil
}

// Function to check if user has supervisor rights
func FindResponderByUUID(uuid string) (bool, error) {

	result, err := findUserRoleByUUID(uuid)

	if err != nil {
		return false, err
	}

	var Access dbtypes.AccessLevel

	return Access.ToInt(result) >= Access.ToInt("responder"), nil
}

// Function to check if user has supervisor rights
func FindSupervisorByUUID(uuid string) (bool, error) {

	result, err := findUserRoleByUUID(uuid)

	if err != nil {
		return false, err
	}

	var Access dbtypes.AccessLevel

	return Access.ToInt(result) >= Access.ToInt("supervisor"), nil
}

func FindAdminByUUID(uuid string) (bool, error) {

	result, err := findUserRoleByUUID(uuid)

	if err != nil {
		return false, err
	}

	var Access dbtypes.AccessLevel

	return Access.ToInt(result) >= Access.ToInt("admin"), nil
}

// Returns UUIDs of files in a case collection
func FindFilesByCase(caseUUID string) ([]string, error) {

	result, err := FindDocsByFilter("Cases", caseUUID, bson.M{})

	if err != nil {
		return nil, err
	}

	var fileList []string

	for _, doc := range result {
		fileList = append(fileList, doc["uuid"].(string))
	}

	return fileList, nil
}

//returns the job
func FindJobStatusByUUID(jobUUID string) (dbtypes.JobStatus, error) {

	result, err := FindDocByFilter("Jobs", "JobQueue", bson.M{"jobuuid": jobUUID})

	if err != nil {
		return "", err
	}

	var jobFromResult dbtypes.Job

	err = result.Decode(&jobFromResult)

	if err != nil {
		return "", err
	}

	return jobFromResult.Status, nil
}

//search for a job;
//this can be used for getting the UUID
func FindJobByFilter(jobFilter interface{}) (dbtypes.Job, error) {

	result, err := FindDocByFilter("Jobs", "JobQueue", bson.M{"jobuuid": jobFilter})

	if err != nil {
		return dbtypes.Job{}, err
	}

	var jobFromResult dbtypes.Job

	err = result.Decode(&jobFromResult)

	if err != nil {
		return dbtypes.Job{}, err
	}

	return jobFromResult, nil

}

// finds jobs that are available inside of the job queue
func FindAvailableJobs(jobTypes []string) ([]dbtypes.Job, error) {
	var decodedJobResult dbtypes.Job
	//var jobResults map[string][]dbtypes.Job = make(map[string][]dbtypes.Job)

	var jobResults []dbtypes.Job = make([]dbtypes.Job, 0)

	//for each type of job
	for _, jobType := range jobTypes {
		//get the results
		results, err := FindDocsByFilter("Jobs", "JobQueue", bson.M{"jobtype": jobType, "status": dbtypes.Queued})

		//hee-hoo error handling
		//go to the next type since apparently no jobs of that type exist
		if err != nil {
			continue
		}

		//marshal each document into a Job struct and then append it to the array
		for _, jobDoc := range results {

			bsonBytes, err := bson.Marshal(jobDoc)

			if err != nil {
				log.Panicln("INTERNAL SERVER ERROR: UNMARSHALLING JOB BSON FAILED")
			}

			err = bson.Unmarshal(bsonBytes, &decodedJobResult)

			if err != nil {
				log.Panicln("INTERNAL SERVER ERROR: UNMARSHALLING JOB BSON FAILED")
			}

			jobResults = append(jobResults, decodedJobResult)

		}
	}

	return jobResults, nil
}

//returns the types of jobs that are available to work on
func FindJobTypes(dbname string, collection string) ([]bson.M, error) {
	// connect to db
	client, ctx, cancel, err := dbConnect()

	// defer closing db connection
	defer dbClose(client, ctx, cancel)

	if err != nil {
		return nil, err
	}

	// get collection
	coll := client.Database(dbname).Collection(collection)

	findOptions := options.Find().SetProjection(bson.M{"jobtype": 1})

	// run find function and get cursor
	cur, err := coll.Find(ctx, bson.D{{}}, findOptions)

	if err != nil {
		return nil, err
	}

	// create slice to hold documents
	var docList []bson.M

	// iterate through the documents
	for cur.Next(context.Background()) {
		// decode cur into a interface
		var doc bson.M
		err := cur.Decode(&doc)

		if err != nil {
			return nil, err
		}
		// log.Println("[DEBUG] internal result: ", doc)

		// append doc to docList
		docList = append(docList, doc)
	}

	// return list of documents
	return docList, nil
}

//retrieves logs by
func GetCaseLogs(c *gin.Context, caseuuid string) ([]bson.M, error) {

	docs, err := FindDocsByFilter("Logs", "Logs", bson.M{"uuid": caseuuid})

	if err != nil {
		return nil, err
	}

	return docs, nil
}
