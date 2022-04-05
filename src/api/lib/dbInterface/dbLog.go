package dbInterface

import "go.mongodb.org/mongo-driver/mongo"

// MakeFile creates a new File struct.
func MakeCaseLog() (*mongo.InsertOneResult, error) { //uuid string, hashes []string, tags []string, filename string, caseUUID string, fileDir string, uploadDate string, viewAccess string, editAccess string) (*mongo.InsertOneResult, error) {

	/*caseName, err := FindCaseNameByUUID(caseUUID)
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
	*/

	result, err = DbSingleInsert(dbName, dbCollection, NewFile)

	return result, err
}
