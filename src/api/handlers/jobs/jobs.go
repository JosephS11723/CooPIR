package jobs
/*
// handles incoming job requests
import (
	"net/http"

	"github.com/JosephS11723/CooPIR/src/api/handlers/iodb"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/gin-gonic/gin"
)

// GetStatus returns the status of a job
func GetInfo(c *gin.Context) {
	var job iodb.Job
	var err error

	// get job id
	jobUUID := c.Param("uuid")

	// get job from db
	job, err = iodb.GetJobInfo(jobUUID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//TODO: get status from job information

	// return job status
	// TODO: change to job status variable
	c.JSON(200, job)
}

// CreateJob creates a new job from parameters given in the request
func CreateJob(c *gin.Context) {
	// create uuid for job
	uuid := dbInterface.MakeUuid()

	// get job name
	jobName := c.PostForm("jobName")

	// empty field check
	if jobName == "" {
		c.JSON(400, gin.H{"error": "jobName is empty"})
		return
	}

	// get job type
	jobType := c.PostForm("jobType")

	// empty field check
	if jobType == "" {
		c.JSON(400, gin.H{"error": "jobType is empty"})
		return
	}

	// get job parameters
	jobParams := c.PostForm("jobParams")

	// empty field check
	if jobParams == "" {
		c.JSON(400, gin.H{"error": "jobParams is empty"})
		return
	}

	// add job to database
	err := iodb.AddJob(uuid, jobName, jobType, jobParams)

	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to add job to database"})
		return
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
		c.JSON(400, gin.H{"error": "jobTypes is empty"})
		return
	}

	// get list of incomplete jobs from database for each job type
	incompleteJobs, err := iodb.GetIncompleteJobs()

	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to get incomplete jobs"})
		return
	}

	// if not jobs found, return not job to user
	if len(incompleteJobs) == 0 {
		c.JSON(200, gin.H{"uuid": "none"})
		return
	}

	// iterate through job types and find a job that matches one of the given job types
	for _, job := range incompleteJobs {
		for _, jobType := range jobTypes {
			if job.JobType == jobType {
				// mark job as in progress in database
				err = iodb.ModifyJobStatus(job.UUID, "in progress")

				if err != nil {
					c.JSON(400, gin.H{"error": "Failed to mark job as in progress"})
					return
				}

				// return job to user and return
				c.JSON(200, gin.H{"uuid": job.UUID})
				return
			}
		}
	}
}

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