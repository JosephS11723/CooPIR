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
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/JosephS11723/CooPIR/src/jobWorker/config"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/resultTypes"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/seaweed"
)

// JobWorker is the struct that wraps the basic JobWorker methods and variables.
type JobWorker struct {
	// JobQueue is the queue of jobs to be done
	JobQueue chan dbtypes.Job

	// JobResultQueue is the queue of job results to be submitted to the api.
	JobResultQueue chan ResultContainer

	// ReturnChan is the channel to return fileuuid results to
	ReturnChan chan string

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

	// client is the HTTP client used to communicate with the api
	client *http.Client
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
type JobFunc func(*dbtypes.Job, chan ResultContainer, chan string) error

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
	j.JobResultQueue = make(chan ResultContainer, 1)

	// create the return channel
	j.ReturnChan = make(chan string, 1)

	// create the job list
	j.jobList = make(map[string]JobFunc)

	// create the done channel
	j.doneChan = make(chan bool)

	// make the waitgroup
	j.wg = sync.WaitGroup{}

	// set the max workers
	j.MaxWorkers = maxWorkers

	// create the http client
	j.client = &http.Client{
		Jar: http.DefaultClient.Jar,
	}

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
		// register with the api
		go j.register()

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
	// if name field on jobresult is blank, set it to the jobuuid
	if result.Name == "" {
		result.Name = result.JobUUID
	}

	// get params
	params := jobResultToParams(&result)

	// infinite loop
	for {
		// sleep for config.SubmitResultInterval seconds
		time.Sleep(time.Duration(config.SubmitResultInterval) * time.Second)

		var req *http.Request
		var err error
		var resp *http.Response

		if r != nil {
			rr, w := io.Pipe()
			mpw := multipart.NewWriter(w)
			go func() {
				var part io.Writer
				defer w.Close()

				if part, err = mpw.CreateFormFile("file", "examplefilename"); err != nil {
					log.Println(err)
					return
				}
				part = io.MultiWriter(part)
				if _, err = io.Copy(part, r); err != nil {
					log.Println(err)
					return
				}
				if err = mpw.Close(); err != nil {
					log.Println(err)
					return
				}
			}()
			// create the request
			// TODO: add reader parsing if there is a file to upload (change nil to reader)
			resp, err = http.Post("http://"+config.ApiName+":"+config.ApiPort+strings.ReplaceAll(submitResultPath, "{jobuuid}", result.JobUUID)+"?"+params, mpw.FormDataContentType(), rr)
		} else {
			// create the request
			// TODO: add reader parsing if there is a file to upload (change nil to reader)
			req, err = http.NewRequest("POST", "http://"+config.ApiName+":"+config.ApiPort+strings.ReplaceAll(submitResultPath, "{jobuuid}", result.JobUUID)+"?"+params, nil)
			// if there is an error, log it and continue
			if err != nil {
				log.Println("Error creating job result request:", err)
				continue
			}

			// send the request
			resp, err = j.client.Do(req)
		}

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

		// get response and sent to return chan
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Println("Error reading job result response:", err)
			continue
		}

		// unmarshal into map
		var m map[string]interface{}
		err = json.Unmarshal(body, &m)

		if err != nil {
			log.Println("Error unmarshalling job result response:", err)
			continue
		}

		// send the response to the return chan
		j.ReturnChan <- m["fileuuid"].(string)

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

	// log
	log.Println("Loaded job", jobTypeName)
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
				err := j.jobList[job.JobType](&job, j.JobResultQueue, j.ReturnChan)

				// if there is an error, log it and send an error message for the job
				if err != nil {
					// send error message
					log.Println("[ERROR workerRoutine]: Error performing job:", err)

					// send error message
					j.JobResultQueue <- ResultContainer{
						JobResult: JobResult{
							JobUUID:    job.JobUUID,
							ResultType: resultTypes.Error,
							Tags:       []string{},
							Relations:  []string{},
							Name:       err.Error(),
							FileUUID:   "na",
							CaseUUID:   job.CaseUUID,
							Done:       true,
						},
						FileReader: nil,
					}
				} else {
					// send done message with no content
					j.JobResultQueue <- ResultContainer{
						JobResult: JobResult{
							JobUUID:    job.JobUUID,
							ResultType: resultTypes.Done,
							Tags:       []string{},
							Relations:  []string{},
							Name:       "na",
							FileUUID:   "na",
							CaseUUID:   job.CaseUUID,
							Done:       true,
						},
						FileReader: nil,
					}
				}

				// void result
				<-j.ReturnChan
			}

		case <-j.doneChan:
			// decrement the waitgroup
			j.wg.Done()

			// exit
			return
		}
	}
}

func (j *JobWorker) register() {
	// get list of keys from jobList
	keys := make([]string, 0, len(j.jobList))
	for k := range j.jobList {
		keys = append(keys, k)
	}

	// for job in jobtypes
	for _, jobType := range keys {
		// create jobtypes as params
		params := registrationToParams(jobType, "worker1")

		// create a client
		client := &http.Client{}

		for {
			// create request
			req, err := http.NewRequest("POST", "http://"+config.ApiName+":"+config.ApiPort+registerPath+"?"+params, nil)

			// if there is an error, log it and return
			if err != nil {
				log.Println("Error creating register request:", err)
				return
			}

			// make the request
			resp, err := client.Do(req)

			// if there is an error, log it and return
			if err != nil {
				log.Println("Error registering worker:", err)
				return
			}

			// if the response is not ok, log it and return
			if resp.StatusCode != http.StatusOK {
				log.Println("Error registering worker:", resp.Status)
				return
			}
			break
		}
	}
}
