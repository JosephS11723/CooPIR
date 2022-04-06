package dbInterface

import (
	"context"
	"log"
	"time"

	"github.com/JosephS11723/CooPIR/src/api/config"
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
		Username:      "myUserAdmin",
		Password:      "securepassword",
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
