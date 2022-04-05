package dbInterface

import (
	"log"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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
