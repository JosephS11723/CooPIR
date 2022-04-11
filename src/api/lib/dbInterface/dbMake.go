package dbInterface

import (
	"errors"
	"log"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/api/lib/security"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

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

	// Log struct
	case dbtypes.Log:
		if collection != "Logs" {
			log.Panicf("[ERROR] Cannot insert data type %s into Log collection", t)
		} else {

			data := data.(dbtypes.Log)

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
func MakeCase(NewCase dbtypes.Case) (*mongo.InsertOneResult, string, error) {
	// Check if case name already exists
	exists, err := DoesCaseExist(NewCase.Name)
	if exists {
		// If case name exists, return error
		return nil, "", errors.New("case name already exists")
	}

	if err != nil {
		return nil, "", err
	}

	// Set db types
	var dbName string = "Cases"
	var dbCollection string = "CaseMetadata"
	id, err := MakeUuid()

	// could not make uuid
	if err != nil {
		return nil, "", err
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

	return result, id, err
}

// MakeFile creates a new File struct.
func MakeFile(uuid string, hashes []string, tags []string, filename string, caseUUID string, fileDir string, uploadDate string, viewAccess string, editAccess string) (*mongo.InsertOneResult, error) {
	caseName, err := FindCaseNameByUUID(caseUUID)
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
func MakeAccess(target string, caseuuid string, user string, time string) (*mongo.InsertOneResult, error) {
	var err error
	var dbName string = "Cases"
	var dbCollection string = "Log"
	id, err := MakeUuid()

	if err != nil {
		return nil, err
	}

	var result *mongo.InsertOneResult

	var NewAccess = dbtypes.Access{
		UUID:   id,
		Target: target,
		User:   user,
		Time:   time,
	}

	result, err = DbSingleInsert(dbName, dbCollection, NewAccess)

	return result, err
}

// MakeUuid creates a new UUID and checks the database to make sure it doesn't already exist.
func MakeUuid() (string, error) {
	var id string
	var exist bool
	var err error
	var Users []string = findCollections("Users")
	var Cases []string = findCollections("Cases")

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
