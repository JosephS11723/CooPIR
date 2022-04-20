package dbtypes

type Worker struct {
	WorkerUUID string       `json:"workeruuid"`
	Name       string       `json:"name"`
	JobType    string       `json:"jobtype"`
	Status     WorkerStatus `json:"status"`
	JoinTime   int          `json:"jointime"`
}

//this is the status
//of
//the worker struct
type WorkerStatus string

//possible values for WorkerStatus
const (
	Running      WorkerStatus = "running"
	Ready        WorkerStatus = "ready"
	Stopped      WorkerStatus = "stopped"
	Disconnected WorkerStatus = "disconnected"
	WorkerError  WorkerStatus = "workererror"
)

func (s WorkerStatus) String() string {
	var return_val string

	switch s {
	case Running:
		return_val = "running"
	case Ready:
		return_val = "ready"
	case Stopped:
		return_val = "stopped"
	case Disconnected:
		return_val = "disconnected"
	case WorkerError:
		return_val = "workererror"
	}

	return return_val
}

//this is just an intermediate value for serializing requests into
type NewWorker struct {
	WorkerUUID string `json:"workeruuid"`
	Name       string `json:"name"`
	JobType    string `json:"jobtype"`
}
