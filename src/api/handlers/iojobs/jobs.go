package iojobs

import (
	"net/http"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
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

	jobUUID := c.Query("jobuuid")

	if jobUUID == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "no job uuid in the query",
			},
		)
	}

	resultType := c.Query("resulttype")

	if resultType == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "no resulttype in the query",
			},
		)
	}

	done := c.Query("resulttype")

	if done == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "no parameter 'done' in the query",
			},
		)
	}

	name := c.Query("name")

	if name == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "no name in the query",
			},
		)
	}

	tags := c.QueryArray("tags")
	relations := c.QueryArray("relations")

}

/*
// SubmitWork submits a job result to the database.
// the job result is sent in pieces which must be added sequentially
func SubmitWork(c *gin.Context) {
	// get job id
	jobUUID := c.Param("uuid")

	// get status
	status := c.PostForm("jobStatus")

	// empty field check
	if status == "" {
		c.JSON(400, gin.H{"error": "status is empty"})
		return
	} else {
		// if job status is error, set job status to error
		if status == "error" {
			// get error message from formdata
			errorMessage := c.PostForm("errorMessage")

			// empty field check
			if errorMessage == "" {
				c.JSON(400, gin.H{"error": "errorMessage is empty"})
				return
			}

			// add error status to job
			err := iodb.ModifyJobStatus(jobUUID, "error")

			if err != nil {
				c.JSON(400, gin.H{"error": "Failed to mark job as error"})
				return
			}

			// add error message to job
			err = iodb.ModifyJobErrorMessage(jobUUID, errorMessage)
		}
	}

	// get reader for job pieces from multipart
	jobPieces := c.Request.MultipartForm.File["jobPieces"]

	var startPiece string
	var pieceCount int = 0

	// add each piece of the job result to the database by iterating through the pieces given
	for pieceIndex, piece := range jobPieces {
		pieceCount = pieceIndex + 1
		// generate jobPiece uuid
		jobPieceUUID := dbInterface.MakeUuid()

		switch pieceIndex {
		case 0:
			// first job piece
			// add job piece to database
			err := iodb.AddJobPiece(jobPieceUUID, jobUUID, piece)

			// set startPiece to jobPieceUUID
			startPiece = jobPieceUUID

			if err != nil {
				// error 500
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add job piece to database"})
				return
			}

			// add job result to database with first piece as head
		default:
			// everything else
		}

		// add jobPiece to database
		err := iodb.AddJobPiece(jobPieceUUID, jobUUID, piece)

		if err != nil {
			// error 500
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add job piece to database"})
			return
		}
	}

	// empty pieces check
	if pieceCount == 0 {
		// error 400
		c.JSON(http.StatusBadRequest, gin.H{"error": "jobPieces is empty"})
		return
	}

	// add job result to database
	err := iodb.AddJobResult(jobUUID, startPiece, pieceCount)

	if err != nil {
		// error 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add job result to database"})
		return
	}

	// mark job as complete
	err = iodb.ModifyJobStatus(jobUUID, "done")

	if err != nil {
		// error 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add job result to database"})
		return
	}

	// return success (200)
	c.Status(http.StatusOK)
}


// GetResults sends the results of a job as a multipart to the client
func GetResults(c *gin.Context) {}
*/
