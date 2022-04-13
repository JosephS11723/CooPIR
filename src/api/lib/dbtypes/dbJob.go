package dbtypes

type Job struct {
	JobUUID       string   `json:"jobuuid"`
	Arguments     []string `json:"arguments"`
	Name          string   `json:"name"`
	JobType       string   `json:"jobtype"`
	Status        Status   `json:"status"`
	StartTime     int      `json:"starttime"`
	EndTime       int      `json:"endtime"`
	JobResultUUID string   `json:"jobresultuuid"`
}

//this is the status
//of
//the job struct
type Status string

//possible values for Status
const (
	Queued     Status = "queued"
	InProgress Status = "in-progress"
	Finished   Status = "finished"
	JobError   Status = "error"
	Cancelled  Status = "cancelled"
)

type NewJob struct {
	Arguments []string `json:"arguments"`
	Name      string   `json:"name"`
	JobType   string   `json:"jobtype"`
}
