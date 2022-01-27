package dbInterface

import (
	// "io"
	"context"
	"log"

	// _ "sync"
	"time"

	"github.com/JosephS11723/CooPIR/src/api/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// This is a user defined method that returns mongo.Client,
// context.Context, context.CancelFunc and error.
// mongo.Client will be used for further database operation.
// context.Context will be used set deadlines for process.
// context.CancelFunc will be used to cancel context and
// resource associated with it.
func DbConnect() (*mongo.Client, context.Context,
	context.CancelFunc, error) {

	uri := config.DBIP

	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.MongoConnectionTimeout)*time.Second)

	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

// This is a user defined method to close resources.
// This method closes mongoDB connection and cancel context.
func DbClose(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			log.Panicln(err)
		}
	}()
}

// This is a user defined method that accepts
// mongo.Client and context.Context
// This method used to ping the mongoDB, return error if any.
func DbPing(client *mongo.Client, ctx context.Context) error {

	// mongo.Client has Ping to ping mongoDB, deadline of
	// the Ping method will be determined by cxt
	// Ping method return error if any occurred, then
	// the error can be handled.
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	log.Println("connected successfully")
	return nil
}
