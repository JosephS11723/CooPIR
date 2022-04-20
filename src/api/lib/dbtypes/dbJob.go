package dbtypes

type Job struct {
	JobUUID       string    `json:"jobuuid"`
	Arguments     []string  `json:"arguments"`
	Name          string    `json:"name"`
	JobType       string    `json:"jobtype"`
	Status        JobStatus `json:"status"`
	StartTime     int       `json:"starttime"`
	EndTime       int       `json:"endtime"`
	JobResultUUID string    `json:"jobresultuuid"`
}

//this is the status
//of
//the job struct
type JobStatus string

//possible values for Status
const (
	Queued     JobStatus = "queued"
	InProgress JobStatus = "in-progress"
	Finished   JobStatus = "finished"
	JobError   JobStatus = "error"
	Cancelled  JobStatus = "cancelled"
)

func (s JobStatus) String() string {
	var return_val string

	switch s {
	case Queued:
		return_val = "queued"
	case InProgress:
		return_val = "in-progress"
	case Finished:
		return_val = "finished"
	case JobError:
		return_val = "error"
	case Cancelled:
		return_val = "cancelled"
	}

	return return_val
}

type NewJob struct {
	Arguments []string `json:"arguments"`
	Name      string   `json:"name"`
	JobType   string   `json:"jobtype"`
}
