package dbInterface

import (
	// "io"
	"context"
	"errors"
	"log"

	// _ "sync"
	"time"

	"github.com/JosephS11723/CooPIR/src/api/config"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/api/lib/security"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// dbConnect returns a mongoDB client, context, cancel function, and error.
func dbConnect() (*mongo.Client, context.Context, context.CancelFunc, error) {
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

// dbClose closes mongoDB connection and cancel context.
func dbClose(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	// CancelFunc to cancel to context
	defer cancel()

	defer func() {

		if err := client.Disconnect(ctx); err != nil {
			log.Panicln(err)
		}
	}()
}

// DbPing is used to ping the mongoDB, return error if any.
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

// DbSingleInsert inserts a single document into a collection.
func DbSingleInsert(dbname string, collection string, data interface{}) (*mongo.InsertOneResult, error) {
	client, ctx, cancel, err := dbConnect()

	// defer closing db connection
	defer dbClose(client, ctx, cancel)

	if err != nil {
		log.Panicln(err)
	}

	switch t := data.(type) {

	// Access struct case
	case dbtypes.Access:
		if collection != "Log" {
			log.Panicf("[ERROR] Cannot insert data type %s into Log collection", t)
		} else {

			data := data.(dbtypes.Access)

			coll := client.Database(dbname).Collection(collection)

			result, err := coll.InsertOne(ctx, data)

			if err != nil {
				return result, err
			}

			return result, nil
		}

	// Case struct case
	case dbtypes.Case:
		if collection != "CaseMetadata" {
			log.Panicf("[ERROR] Cannot insert data type %s into Case collection", t)
		} else {

			data := data.(dbtypes.Case)

			coll := client.Database(dbname).Collection(collection)

			result, err := coll.InsertOne(ctx, data)

			if err != nil {
				return result, err
			}

			return result, nil
		}

	// File struct case
	// TODO: Replace Sanity check with a more robust check
	case dbtypes.File:
		// if collection != "File" {
		// 	log.Panicf("[ERROR] Cannot insert data type %s into File collection", t)
		// } else {

		data := data.(dbtypes.File)

		coll := client.Database(dbname).Collection(collection)

		result, err := coll.InsertOne(ctx, data)

		if err != nil {
			return result, err
		}

		return result, nil
		// }

	// User struct case
	case dbtypes.User:
		if collection != "UserMetadata" {
			log.Panicf("[ERROR] Cannot insert data type %s into User collection", t)
		} else {

			data := data.(dbtypes.User)

			coll := client.Database(dbname).Collection(collection)

			result, err := coll.InsertOne(ctx, data)

			if err != nil {
				return result, err
			}

			return result, nil
		}

	// default case: panic
	default:
		log.Panic("[ERROR] Unknown type for db insert!")
	}

	return nil, nil
}

// MakeUser creates a new User struct.
//func MakeUser(name string, email string, role string, cases []string, password string) (*mongo.InsertOneResult, error) {
func MakeUser(user dbtypes.NewUser) (*mongo.InsertOneResult, error) {
	// check if user email already exists
	exists, err := doesEmailExist(user.Email)
	if exists {
		// if email exists, return error
		return nil, errors.New("email already exists")
	}

	if err != nil {
		return nil, err
	}

	// Hash password
	saltedHash, err := security.HashPass(user.Password)
	if err != nil {
		return nil, err
	}

	// Set db types
	var dbName string = "Users"
	var dbCollection string = "UserMetadata"

	// Make unique id
	id, err := MakeUuid()

	// could not make uuid
	if err != nil {
		return nil, err
	}

	// Set user struct
	var NewUser = dbtypes.User{
		UUID:       id,
		Name:       user.Name,
		Email:      user.Email,
		Role:       user.Role,
		Cases:      user.Cases,
		SaltedHash: saltedHash,
	}

	result, err := DbSingleInsert(dbName, dbCollection, NewUser)

	return result, err
}

// MakeCase creates a new Case struct.
//func MakeCase(NewCase dbttypes.Case) *mongo.InsertOneResult {
func MakeCase(NewCase dbtypes.Case) (*mongo.InsertOneResult, error) {
	// Check if case name already exists
	exists, err := DoesCaseExist(NewCase.Name)
	if exists {
		// If case name exists, return error
		return nil, errors.New("case name already exists")
	}

	if err != nil {
		return nil, err
	}

	// Set db types
	var dbName string = "Cases"
	var dbCollection string = "CaseMetadata"
	id, err := MakeUuid()

	// could not make uuid
	if err != nil {
		return nil, err
	}

	NewCase.UUID = id
	var result *mongo.InsertOneResult

	// Replaced
	/*
		var NewCase = dbtypes.Case{
			UUID:          id,
			Name:          name,
			Date_created:  dateCreated,
			View_access:   viewAccess,
			Edit_access:   editAccess,
			Collaborators: collaborators,
		}
	*/

	client, ctx, cancel, err := dbConnect()

	// defer closing db connection
	defer dbClose(client, ctx, cancel)

	if err != nil {
		log.Panicln(err)
	}

	client.Database("Cases").CreateCollection(ctx, id)

	result, err = DbSingleInsert(dbName, dbCollection, NewCase)

	return result, err
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

	var userEmail string = dbUser.Email

	return userEmail, nil
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

	var userUUID string = dbUser.UUID

	return userUUID, err
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

	var caseName string = dbCase.Name

	return caseName, err
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

	var caseUUID string = dbCase.UUID

	return caseUUID, err
}

// MakeFile creates a new File struct.
func MakeFile(uuid string, hashes []string, tags []string, filename string, caseName string, fileDir string, uploadDate string, viewAccess string, editAccess string) (*mongo.InsertOneResult, error) {
	caseUUID, err := FindCaseUUIDByName(caseName)
	if err != nil {
		return nil, err
	}

	var dbName string = "Cases"
	var dbCollection string = caseUUID
	var result *mongo.InsertOneResult

	var NewFile = dbtypes.File{
		UUID:        uuid,
		MD5:         hashes[0],
		SHA1:        hashes[1],
		SHA256:      hashes[2],
		SHA512:      hashes[3],
		Tags:        tags,
		Filename:    filename,
		Case:        caseName,
		File_dir:    fileDir,
		Upload_date: uploadDate,
		View_access: viewAccess,
		Edit_access: editAccess,
	}

	result, err = DbSingleInsert(dbName, dbCollection, NewFile)

	return result, err
}

// MakeAccess creates a new Access struct.
func MakeAccess(filename string, user string, date string) (*mongo.InsertOneResult, error) {
	var dbName string = "Cases"
	var dbCollection string = "Log"
	var result *mongo.InsertOneResult
	
	id, err := MakeUuid()
	if err != nil {
		return nil, err
	}

	var NewAccess = dbtypes.Access{
		UUID:     id,
		Filename: filename,
		User:     user,
		Date:     date,
	}

	result, err = DbSingleInsert(dbName, dbCollection, NewAccess)

	return result, err
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

// TODO: Logic is broken in the check against the database. FIX
// MakeUuid creates a new UUID and checks the database to make sure it doesn't already exist.
func MakeUuid() (string, error) {
	var id string
	var exist bool
	var err error
	var Users []string = []string{"UserMetadata"}
	var Cases []string = []string{"CaseMetadata", "File", "Log"}
	// TODO: check more collections (all the collections for all the cases)

	// Loop that makes a uuid and checks if it already exists in the database.
	// keep looping until it doesn't exist.
	for {
		id = uuid.New().String()

		for _, collection := range Users {
			exist, err = doesUuidExist("Users", collection, id)

			if err != nil {
				return "", err
			}

			if exist {
				break
			}
		}

		if exist {
			continue
		}

		for _, collection := range Cases {
			exist, err = doesUuidExist("Cases", collection, id)

			if err != nil {
				return "", err
			}

			if exist {
				break
			}
		}

		if !exist {
			break
		}

	}

	return id, nil
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

// updateDoc modifies a single document's information in the database.
func UpdateDoc(dbName string, dbCollection string, filter bson.M, updates bson.D) *mongo.UpdateResult {
	// connect to db
	client, ctx, cancel, err := dbConnect()

	// defer closing db connection
	defer dbClose(client, ctx, cancel)

	if err != nil {
		log.Panicln(err)
	}

	// get collection
	coll := client.Database(dbName).Collection(dbCollection)

	result, err := coll.UpdateOne(ctx, filter, updates)

	if err != nil {
		log.Panicln(err)
	}

	return result
}

//wrapper around the UpdateDoc function specifically for updating cases
func UpdateCase(dbName string, dbCollection string, caseUpdate dbtypes.UpdateDoc) *mongo.UpdateResult {
	//get the filter, which will act as a bson.M
	var filter map[string]interface{} = caseUpdate.Filter

	//get the update field and the proceed with removing UUID and Date_created
	var unchecked_update map[string]interface{} = caseUpdate.Update

	delete(unchecked_update, "uuid")

	delete(unchecked_update, "dateCreated")

	//this constructs the update bson.D
	var update bson.D = bson.D{{"$set", unchecked_update}}

	return UpdateDoc(dbName, dbCollection, filter, update)
}

// wrapper around the UpdateDoc function specifically for updating cases
func UpdateUser(dbName string, dbCollection string, caseUpdate dbtypes.UpdateDoc) *mongo.UpdateResult {
	// get the filter, which will act as a bson.M
	var filter map[string]interface{} = caseUpdate.Filter

	// get the update field and the proceed with removing UUID and Date_created
	var unchecked_update map[string]interface{} = caseUpdate.Update

	delete(unchecked_update, "uuid")

	delete(unchecked_update, "saltedhash")

	delete(unchecked_update, "role")

	// this constructs the update bson.D
	var update bson.D = bson.D{{"$set", unchecked_update}}

	return UpdateDoc(dbName, dbCollection, filter, update)
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
