package dbInterface

import (
	// "io"
	"context"
	"log"

	// _ "sync"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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

// FindDocsByFilter finds multiple documents in a collection by a filter and RETURN a slice of documents (bson.m)
func FindDocsByFilter(dbname string, collection string, filter bson.M) []bson.M {
	// connect to db
	client, ctx, cancel, err := dbConnect()

	// defer closing db connection
	defer dbClose(client, ctx, cancel)

	if err != nil {
		log.Panicln(err)
	}

	// get collection
	coll := client.Database(dbname).Collection(collection)

	// run find function and get cursor
	cur, err := coll.Find(ctx, filter)

	if err != nil {
		log.Panicln(err)
	}

	// create slice to hold documents
	var docList []bson.M

	// iterate through the documents
	for cur.Next(context.Background()) {
		// decode cur into a interface
		var doc bson.M
		err := cur.Decode(&doc)

		// remove _id from doc
		delete(doc, "_id")

		if err != nil {
			log.Panicln(err)
		}
		log.Println("[DEBUG] internal result: ", doc)

		// append doc to docList
		docList = append(docList, doc)
	}

	// return list of documents
	return docList
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
