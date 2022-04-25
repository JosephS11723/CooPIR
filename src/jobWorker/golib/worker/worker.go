package worker

// import when library is public
/*import (
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
)*/

// temporary local import
import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/JosephS11723/CooPIR/src/jobWorker/config"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/seaweed"
)

// JobWorker is the struct that wraps the basic JobWorker methods and variables.
type JobWorker struct {
	// JobQueue is the queue of jobs to be done
	JobQueue chan dbtypes.Job

	// JobResultQueue is the queue of job results to be submitted to the api.
	JobResultQueue chan ResultContainer

	// MaxWorkers is the maximum number of workers routines to be spawned
	MaxWorkers int `json:"maxworkers default:2"`

	// curJobs is the current number of in-progress jobs
	curJobs int `default:"0"`

	// curJobsLock is the mutex for curJobs
	curJobsLock sync.Mutex

	// jobList is the map of jobstrings that can be performed by this worker by an associated function
	jobList map[string]JobFunc

	// wg is the waitgroup for all the worker routines
	wg sync.WaitGroup

	// done is a bool that is set to true when the worker is set to exit
	done bool `default:"false"`

	// doneChan is a channel that allows blocked workers to unblock obtaining a job
	doneChan chan bool
}

// JobResult is the struct that wraps the information sent to the api
type JobResult struct {
	// JobUUID is the uuid of the job
	JobUUID string `json:"jobuuid"`

	// resultType is the type of result
	ResultType string `json:"resulttype"`

	// tags are the tags associated with the result
	Tags []string `json:"tags"`

	// Relations are the relations associated with the result
	Relations []string `json:"relations"`

	// name is the name of the result
	Name string `json:"name"`

	// done is a bool that indicates if the result is the final result
	Done bool `json:"done"`

	// FileUUID references the file that the result is associated with
	FileUUID string `json:"fileuuid" default:"new"`

	// caseUUID references the case that the result is associated with
	CaseUUID string `json:"caseuuid"`
}

var (
	// ErrWorkerDone is the error returned when the worker is done
	ErrWorkerDone = errors.New("worker is done")
)

// JobFunc is the function type that is called when a job is performed
type JobFunc func(*dbtypes.Job, chan ResultContainer)

// ResultContainer contains the JobResult and the io.reader for any files to be uploaded
type ResultContainer struct {
	JobResult  JobResult
	FileReader io.Reader
}

// NewJobWorker creates a new JobWorker
func NewJobWorker(maxWorkers int) *JobWorker {
	// create a new JobWorker
	j := JobWorker{}

	// create the job queue
	j.JobQueue = make(chan dbtypes.Job)

	// create the job result queue
	j.JobResultQueue = make(chan ResultContainer, 100)

	// create the job list
	j.jobList = make(map[string]JobFunc)

	// create the done channel
	j.doneChan = make(chan bool)

	// make the waitgroup
	j.wg = sync.WaitGroup{}

	// set the max workers
	j.MaxWorkers = maxWorkers

	// return the new JobWorker
	return &j
}

// getWorkLoop is a routine that gets jobs from the api and queues them
func (j *JobWorker) getWorkLoop() {
	for {
		// if curJobs is less than max workers
		if j.curJobs < j.MaxWorkers {
			// ask server for work
			job, err := j.GetJob()

			// if there is an error, log it and continue
			if err != nil {
				log.Println("Error getting job:", err)
				continue
			}

			// increment curJobs
			j.curJobsLock.Lock()
			j.curJobs++
			j.curJobsLock.Unlock()

			// add the job to the queue
			j.JobQueue <- job
		}
	}
}

// submitWorkLoop is a routine that submits job results to the api
func (j *JobWorker) submitWorkLoop() {
	for {
		// wait for a job result
		jobContainer := <-j.JobResultQueue

		// submit the job result
		err := j.SubmitJob(jobContainer.JobResult, jobContainer.FileReader)

		// if there is an error, log it and continue
		if err != nil {
			log.Println("Error submitting job result:", err)
			// add job back into the finished queue
			j.JobResultQueue <- jobContainer
			continue
		}

		// decrement curJobs
		j.curJobsLock.Lock()
		j.curJobs--
		j.curJobsLock.Unlock()
	}
}

// Start starts the JobWorker
func (j *JobWorker) Start() {
	go func() {
		// mount the filer
		go seaweed.MountAllFiles()

		// spawn workers
		for i := 0; i < j.MaxWorkers; i++ {
			// add a worker to the waitgroup
			j.wg.Add(1)

			// spawn a worker
			go j.workerRoutine()
		}

		// spawn the get work loop
		go j.getWorkLoop()

		// spawn the submit work loop
		go j.submitWorkLoop()
	}()
}

// Stop stops the JobWorker.
func (j *JobWorker) Stop() error {
	// set the done bool to true for all the workers
	j.done = true

	// write to the done channel to unblock the workers
	for i := 0; i < j.MaxWorkers; i++ {
		j.doneChan <- true
	}

	// wait for the routines and the waitgroup to finish
	j.wg.Wait()

	return nil
}

// SubmitJob submits a job result to the api
func (j *JobWorker) SubmitJob(result JobResult, r io.Reader) error {
	// create a new http client
	client := &http.Client{}

	// get params
	params := jobResultToParams(&result)

	// infinite loop
	for {
		// sleep for config.SubmitResultInterval seconds
		time.Sleep(time.Duration(config.SubmitResultInterval) * time.Second)

		// create the request
		// TODO: add reader parsing if there is a file to upload (change nil to reader)
		req, err := http.NewRequest("POST", "http://"+config.ApiName+":"+config.ApiPort+strings.ReplaceAll(submitResultPath, "{jobuuid}", result.JobUUID)+"?"+params, r)

		// if there is an error, log it and continue
		if err != nil {
			log.Println("Error creating job result request:", err)
			continue
		}

		// send the request
		resp, err := client.Do(req)

		// if there is an error, log it and continue
		if err != nil {
			log.Println("Error sending job result request:", err)
			continue
		}

		// if the response is not 200, log it and continue
		if resp.StatusCode != 200 {
			log.Println("Error submitting job result status:", resp.Status)
			continue
		}

		// break
		break
	}

	return nil
}

// GetJob gets a job from the api
func (j *JobWorker) GetJob() (dbtypes.Job, error) {
	// create a client
	client := &http.Client{}

	// get list of keys from jobList
	keys := make([]string, 0, len(j.jobList))
	for k := range j.jobList {
		keys = append(keys, k)
	}

	// create jobtypes as params
	params := jobTypesToParams(keys)

	// unmarshal the response body
	var job dbtypes.Job

	// fetch loop
	for {
		// sleep for config.GetJobInterval seconds
		time.Sleep(time.Duration(config.GetJobInterval) * time.Second)

		// if done, break loop and return error
		if j.done {
			return dbtypes.Job{}, ErrWorkerDone
		}

		// get job from api
		// create get request for job
		req, err := http.NewRequest("GET", "http://"+config.ApiName+":"+config.ApiPort+getWorkPath+"?"+params, nil)

		// if there is an error, log it and continue
		if err != nil {
			log.Println("Error creating get job request:", err)
			continue
		}

		// make the request
		resp, err := client.Do(req)

		// if there is an error, log it and continue
		if err != nil {
			log.Println("Error getting job:", err)
			continue
		}

		// if the response is not ok, log it and continue
		if resp.StatusCode != http.StatusOK {
			log.Println("Error getting job:", resp.Status)
			continue
		}

		// read the response body
		body, err := ioutil.ReadAll(resp.Body)

		// if there is an error, log it and continue
		if err != nil {
			log.Println("Error reading job response body:", err)
			continue
		}

		// unmarshal the response body
		err = json.Unmarshal(body, &job)

		// if there is an error, log it and continue
		if err != nil {
			log.Println("Error unmarshalling job response body:", err)
			continue
		}

		// return the job
		return job, nil
	}
}

// AddJobWithFunction adds a function to the list of possible jobs this worker can perform
func (j *JobWorker) AddJobWithFunction(jobTypeName string, jobFunction JobFunc) {
	// add job to job list
	j.jobList[jobTypeName] = jobFunction
}

// workerRoutine is the main routine spawned by the worker struct that creates the "workers" who take in jobs from the queue and perform them
func (j *JobWorker) workerRoutine() {
	// infinite loop
	for {
		// select between getting a job or exiting
		select {
		case job := <-j.JobQueue:
			// print job struct for debug
			fmt.Printf("Received Job: %+v\n", job)
			// if the job is not in the list of jobs this worker can perform, log it and continue
			if _, ok := j.jobList[job.JobType]; !ok {
				log.Println("[ERROR workerRoutine]: Job type not supported:", job.JobType)

				// reduce curjobs
				j.curJobsLock.Lock()
				j.curJobs--
				j.curJobsLock.Unlock()

				continue
			} else {
				// perform the job
				j.jobList[job.JobType](&job, j.JobResultQueue)
			}

		case <-j.doneChan:
			// decrement the waitgroup
			j.wg.Done()

			// exit
			return
		}
	}
}
