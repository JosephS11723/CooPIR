package iojobs

import (
	"log"
	"net/http"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/api/lib/seaweedInterface"
	"github.com/gin-gonic/gin"
)

// handles incoming job requests

// GetStatus returns the status of a job
func GetStatus(c *gin.Context) {

	var status dbtypes.JobStatus
	var err error

	// get job id
	jobUUID := c.Query("uuid")

	if jobUUID == "" {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": "query did not contain job uuid",
			},
		)
	}

	// get job from db
	status, err = dbInterface.FindJobStatusByUUID(jobUUID)

	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, status.String())
}

//just gets the job document
func GetJobInfo(c *gin.Context) {

	var job dbtypes.Job
	var json_request map[string]interface{}
	var err error

	err = c.BindJSON(&json_request)

	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": "could not serialize request info",
			},
		)
	}

	// get job from db
	job, err = dbInterface.FindJobByFilter(json_request)

	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, job)
}

// CreateJob creates a new job from parameters given in the request
func CreateJob(c *gin.Context) {

	var new_job_request dbtypes.NewJob

	err := c.ShouldBind(&new_job_request)

	if err != nil {

		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error":  "request could not be marshalled into a NewJob",
				"Andrew": "of the Merrow sort",
			},
		)

	}

	// add job to database
	uuid, err := dbInterface.MakeJob(new_job_request)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Failed to add job to database"})
	}

	// return job id
	c.JSON(200, gin.H{"uuid": uuid})
}

// GetWork returns a job that matches one of the given capable job types
func GetWork(c *gin.Context) {
	// get job type
	jobTypes := c.QueryArray("jobTypes")

	// empty field check
	if len(jobTypes) == 0 {
		c.JSON(400, gin.H{"error": "jobtype is empty"})
		return
	}

	// get list of incomplete jobs from database for each job type
	incompleteJobs, err := dbInterface.FindAvailableJobs(jobTypes)

	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to get incomplete jobs"})
		return
	}

	// if not jobs found, return not job to user
	if len(incompleteJobs) == 0 {
		c.JSON(404, gin.H{"uuid": "none"})
		return
	}

	// iterate through job types and find a job that matches one of the given job types
	//I literally hate this nested for loop nonsense, but this is how golang works
	for jobListType, jobList := range incompleteJobs {

		//check to see if there is a job type that the worker can do
		for _, jobType := range jobTypes {
			if jobType == jobListType {

				//get the first out of the list;
				//it should be in order because of indexing in collection
				job_to_send := jobList[0]

				// mark job as in progress in database
				err = dbInterface.ModifyJobStatus(job_to_send.JobUUID, dbtypes.InProgress)

				//if err, then err
				if err != nil {
					c.JSON(400, gin.H{"error": "Failed to mark job as in progress"})
					return
				}

				// return job to user and return
				c.JSON(200, gin.H{"uuid": job_to_send.JobUUID})
				return

			}
		}
	}
}

//this is just for receiving work from workers
func SubmitWork(c *gin.Context) {

	var fileUUID string

	//the next bunch of code is basically just looking for the query parameters
	jobUUID := c.Query("jobuuid")

	if jobUUID == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "no job uuid in the query",
			},
		)
		return
	}

	caseUUID := c.Query("caseuuid")

	if jobUUID == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "no case uuid in the query",
			},
		)
		return
	}

	resultType := c.Query("resulttype")

	if resultType == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "no resulttype in the query",
			},
		)
		return
	}

	done := c.Query("done")

	if done == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "no parameter 'done' in the query",
			},
		)
		return
	}

	name := c.Query("name")

	if name == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "no name in the query",
			},
		)
		return
	}

	//tags and relations can be empty arrays, so no check
	tags := c.QueryArray("tags")
	relations := c.QueryArray("relations")

	//find the status of the job; if not valid status for modifying, return error
	status, err := dbInterface.FindJobStatusByUUID(jobUUID)

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{"error": "could not find job by uuid"},
		)
		return
	}

	if status != dbtypes.InProgress {

		c.JSON(
			http.StatusConflict,
			gin.H{"error": "target job is not In-Progess and cannot be modified by worker"},
		)
		return
	}

	//the only result types are 'modifyFile' and 'createFile'
	switch resultType {

	//since this is modifying a particular existing file, then get the fileuuid query param
	case "modifyFile":

		fileUUID = c.Query("fileuuid")

		if fileUUID == "" {

			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "no fileuuid provided"},
			)
			return

		}

		//call this function to try to modify the tags and relations for the file
		err := dbInterface.ModifyJobTagsAndRelations(fileUUID, caseUUID, tags, relations)

		if err != nil {

			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()},
			)
			return
		}

	//if creating a file, then try to create the file--if there is an error from Seaweed,
	//then that file already exists (surprisingly or unsurprisingly), so then just modify
	//the tags and relations
	case "createFile":

		filestream, _, err := c.Request.FormFile("file")

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no file received"})
			log.Panicln(err)
		}

		fileUUID, err = seaweedInterface.POSTFileJob(caseUUID, filestream)

		if err != nil {

			if err.Error() == "file already exists" {

				err = dbInterface.ModifyJobTagsAndRelations(fileUUID, caseUUID, tags, relations)

				if err != nil {
					c.JSON(
						http.StatusInternalServerError,
						gin.H{"error": err.Error()},
					)
					return
				}

			} else {

				c.JSON(
					http.StatusInternalServerError,
					gin.H{"error": err.Error()},
				)
				return
			}

		}

	default:
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "resulttype not valid"},
		)
		return
	}

	if done == "true" {

		err = dbInterface.ModifyJobStatus(jobUUID, dbtypes.InProgress)

		if err != nil {
			log.Panicln("ERROR: could not modify job status from current to 'Finished'")
		}

	}

	c.JSON(http.StatusOK, gin.H{"jobuuid": jobUUID, "fileuuid": fileUUID})
}

// GetResults sends the results of a job as a multipart to the client
func GetResults(c *gin.Context) {}
