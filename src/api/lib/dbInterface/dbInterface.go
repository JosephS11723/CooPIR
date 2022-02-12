package dbInterface

import (
	// "io"
	"context"
	"log"

	// _ "sync"
	"time"


	"github.com/JosephS11723/CooPIR/src/api/config"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/api/lib/security"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func dbConnect() (*mongo.Client, context.Context,
	context.CancelFunc, error) {

	var uri string = config.DBIP

	// Set client options
	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-256",
		AuthSource:    "admin",
		Username:      "api",
		Password:      "bingusmalingus",
	}
	clientOpts := options.Client().ApplyURI(uri).
		SetAuth(credential)

	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.MongoConnectionTimeout)*time.Second)

	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, clientOpts)
	return client, ctx, cancel, err
}

// This method closes mongoDB connection and cancel context.
func dbClose(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	// CancelFunc to cancel to context
	defer cancel()

	defer func() {

		if err := client.Disconnect(ctx); err != nil {
			log.Panicln(err)
		}
	}()
}

// This method used to ping the mongoDB, return error if any.
func DbPing() error {

	client, ctx, cancel, err := dbConnect()
	if err != nil {
		log.Panicln(err)
	}

	defer dbClose(client, ctx, cancel)

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	log.Println("connected successfully")
	return nil
}

func DbSingleInsert(dbname string, collection string, data interface{}) *mongo.InsertOneResult {

	client, ctx, cancel, err := dbConnect()
	if err != nil {
		log.Panicln(err)
	}

	defer dbClose(client, ctx, cancel)

	switch t := data.(type) {

	case dbtypes.Access:
		if collection != "Log" {
			log.Panicf("[ERROR] Cannot insert data type %s into Log collection", t)
		} else {

			data := data.(dbtypes.Access)

			coll := client.Database(dbname).Collection(collection)

			result, err := coll.InsertOne(ctx, data)

			if err != nil {
				log.Panicln(err)
			}

			return result
		}

	case dbtypes.Case:
		if collection != "Case" {
			log.Panicf("[ERROR] Cannot insert data type %s into Case collection", t)
		} else {

			data := data.(dbtypes.Case)

			coll := client.Database(dbname).Collection(collection)

			result, err := coll.InsertOne(ctx, data)

			if err != nil {
				log.Panicln(err)
			}

			return result
		}

	case dbtypes.File:
		if collection != "File" {
			log.Panicf("[ERROR] Cannot insert data type %s into File collection", t)
		} else {

			data := data.(dbtypes.File)

			coll := client.Database(dbname).Collection(collection)

			result, err := coll.InsertOne(ctx, data)

			if err != nil {
				log.Panicln(err)
			}

			return result
		}

	case dbtypes.User:
		if collection != "User" {
			log.Panicf("[ERROR] Cannot insert data type %s into User collection", t)
		} else {

			data := data.(dbtypes.User)

			coll := client.Database(dbname).Collection(collection)

			result, err := coll.InsertOne(ctx, data)

			if err != nil {
				log.Panicln(err)
			}

			return result
		}

	default:
		log.Panic("[ERROR] Unknown type for db intsert!")
	}

	return nil
}

func MakeUser(uuid string, name string, email string, role string, cases []string, password string) dbtypes.User {

	var saltedhash string = security.HashPass(password)

	var NewUser = dbtypes.User{
		UUID:       uuid,
		Name:       name,
		Email:      email,
		Role:       role,
		Cases:      cases,
		SaltedHash: saltedhash,
	}

	return NewUser
}

func MakeCase(uuid string, name string, dateCreated string, viewAccess string, editAccess string, collaborators []string) dbtypes.Case {

	var NewCase = dbtypes.Case{
		UUID:          uuid,
		Name:          name,
		Date_created:  dateCreated,
		View_access:   viewAccess,
		Edit_access:   editAccess,
		Collaborators: collaborators,
	}

	return NewCase
}

func MakeFile(uuid string, hash string, filename string, caseName string, fileDir string, uploadDate string, viewAccess string, editAccess string) dbtypes.File {

	var NewFile = dbtypes.File{
		UUID:        uuid,
		Hash:        hash,
		Filename:    filename,
		Case:        caseName,
		File_dir:    fileDir,
		Upload_date: uploadDate,
		View_access: viewAccess,
		Edit_access: editAccess,
	}

	return NewFile
}

func MakeAccess(uuid string, filename string, user string, date string) dbtypes.Access {

	var NewAccess = dbtypes.Access{
		UUID:     uuid,
		Filename: filename,
		User:     user,
		Date:     date,
	}

	return NewAccess

}

// Find multiple documents in a collection by a filter and RETURN a slice of documents (bson.m)
func FindDocsByFilter(dbname string, collection string, filter bson.M) []bson.M {

	client, ctx, cancel, err := dbConnect()
	if err != nil {
		log.Panicln(err)
	}

	defer dbClose(client, ctx, cancel)

	coll := client.Database(dbname).Collection(collection)

	cur, err := coll.Find(ctx, filter)

	if err != nil {
		log.Panicln(err)
	}

	var docList []bson.M

	for cur.Next(context.Background()) {
		// Decode cur into a interface
		var doc bson.M
		err := cur.Decode(&doc)

		// remove _id from doc
		delete(doc, "_id")

		if err != nil {
			log.Panicln(err)
		}
		log.Println("[DEBUG] internal result: ", doc)
		docList = append(docList, doc)
	}

	return docList
}

// Find a single document in a collection by a filter and RETURN a document (*mongo.SingleResult)
func FindDocByFilter(dbname string, collection string, filter bson.M) *mongo.SingleResult {

	client, ctx, cancel, err := dbConnect()
	if err != nil {
		log.Panicln(err)
	}

	defer dbClose(client, ctx, cancel)

	// Get the collection
	coll := client.Database(dbname).Collection(collection)

	// Find the document
	var result *mongo.SingleResult = coll.FindOne(ctx, filter)

	// If error, panic
	if result.Err() != nil {
		log.Panicln(result.Err())
	}

	// Return the result
	return result
}

// Check if UUID exists in the collection
func DoesUuidExist(dbname string, collection string, uuid string) bool {

	client, ctx, cancel, err := dbConnect()
	if err != nil {
		log.Panicln(err)
	}

	defer dbClose(client, ctx, cancel)

	// Get the collection
	coll := client.Database(dbname).Collection(collection)

	// Find the document
	var result *mongo.SingleResult = coll.FindOne(ctx, bson.M{"uuid": uuid})

	// Return true if document exists
	return result.Err() != mongo.ErrNoDocuments

}
