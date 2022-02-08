package dbInterface

import (
	// "io"
	"context"
	"log"

	// _ "sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/JosephS11723/CooPIR/src/api/config"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func DbConnect() (*mongo.Client, context.Context,
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
func DbClose(client *mongo.Client, ctx context.Context,
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
func DbPing(client *mongo.Client, ctx context.Context) error {

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	log.Println("connected successfully")
	return nil
}

func DbSingleInsert(client *mongo.Client, ctx context.Context, dbname string,
	collection string, data interface{}) *mongo.InsertOneResult {

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

func PassHash(password string) string {

	// Generate a salted hash of the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panicln(err)
	}

	return string(hash)
}

func MakeUser(name string, email string, role string, cases []string, password string) dbtypes.User {

	var saltedhash string = PassHash(password)

	var NewUser = dbtypes.User{
		Name:       name,
		Email:      email,
		Role:       role,
		Cases:      cases,
		SaltedHash: saltedhash,
	}

	return NewUser
}

func MakeCase(name string, dateCreated string, viewAccess string, editAccess string, collaborators []string) dbtypes.Case {

	var NewCase = dbtypes.Case{
		Name:          name,
		Date_created:  dateCreated,
		View_access:   viewAccess,
		Edit_access:   editAccess,
		Collaborators: collaborators,
	}

	return NewCase
}

func MakeFile(hash string, filename string, caseName string, fileDir string, uploadDate string, viewAccess string, editAccess string) dbtypes.File {

	var NewFile = dbtypes.File{
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

func MakeAccess(filename string, user string, date string) dbtypes.Access {

	var NewAccess = dbtypes.Access{
		Filename: filename,
		User:     user,
		Date:     date,
	}

	return NewAccess

}

func FindDocsByFilter(client *mongo.Client, ctx context.Context, dbname string, collection string, filter bson.M) []string {

	coll := client.Database(dbname).Collection(collection)

	cur, err := coll.Find(ctx, filter)

	if err != nil {
		log.Panicln(err)
	}

	var docList []string

	for cur.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()
		result := struct {
			Foo string
			Bar int32
		}{}
		err := cur.Decode(&result)
		if err != nil {
			log.Panicln(err)
		}
		docList = append(docList, result.Foo)
	}

	return docList
}
