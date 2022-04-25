package dbInterface

import (
	"fmt"
	"log"

	"github.com/JosephS11723/CooPIR/src/api/lib/coopirutil"
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

	//get the update field and the proceed with removing UUID and dateCreated
	var unchecked_update map[string]interface{} = caseUpdate.Update

	// TODO: make sure user has permission to mess with case information. also check for used name so no name duplicates or collisions

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

	// get the update field and the proceed with removing UUID and dateCreated
	var unchecked_update map[string]interface{} = caseUpdate.Update

	delete(unchecked_update, "uuid")

	delete(unchecked_update, "saltedhash")

	delete(unchecked_update, "role")

	// this constructs the update bson.D
	var update bson.D = bson.D{{"$set", unchecked_update}}

	return UpdateDoc(dbName, dbCollection, filter, update)
}

//this just modifies the job document in the Jobs server
func ModifyJobStatus(jobUUID string, status dbtypes.JobStatus) error {

	_ = UpdateDoc("Jobs", "JobQueue", bson.M{"jobuuid": jobUUID}, bson.D{{"$set", bson.M{"status": status}}})

	return nil
}

//find the file by UUID in a case, then append to that file's tags and relations
func ModifyJobTagsAndRelations(fileUUID string, caseUUID string, tags []string, relations []string) error {
	//go through the ceremony of finding the file, decoding it, and error-checking
	result, err := FindDocByFilter("Cases", caseUUID, bson.M{"uuid": fileUUID})

	if err != nil {
		return fmt.Errorf("could not find file %s in case %s", fileUUID, caseUUID)
	}

	var file dbtypes.File

	err = result.Decode(&file)

	if err != nil {
		return fmt.Errorf("could not decode file %s from mongo result", fileUUID)
	}

	//iterate through (since Go doesn't have a simple way of doing this >:( ) and check
	//each tag to see if it already exists
	/*
		for _, tag := range tags {

			if len(file.Tags) > 0 {

				for _, exisiting_tag := range file.Tags {

					if tag != exisiting_tag {

						file.Tags = append(file.Tags, tag)

					}

				}
			} else {
				file.Tags = append(file.Tags, tag)

			}

		}
	*/

	for _, tag := range tags {

		file.Tags = append(file.Tags, tag)

	}

	file.Tags = coopirutil.RemoveDuplicateStr(file.Tags)

	file.Relations = coopirutil.RemoveDuplicateStr(file.Relations)

	for _, relation := range relations {

		file.Relations = append(file.Relations, relation)

	}

	file.Relations = coopirutil.RemoveDuplicateStr(file.Relations)

	/*
		//same thing but for relations
		for _, relation := range relations {

			if len(file.Relations) > 0 {

				for _, exisiting_relation := range file.Relations {

					if relation != exisiting_relation {

						file.Relations = append(file.Relations, relation)

					}

				}
			} else {
				file.Relations = append(file.Relations, relation)

			}
		}*/

	UpdateDoc(
		"Cases",
		caseUUID,
		bson.M{"uuid": fileUUID},
		bson.D{{"$set", bson.M{"tags": file.Tags, "relations": file.Relations}}},
	)

	return nil
}
