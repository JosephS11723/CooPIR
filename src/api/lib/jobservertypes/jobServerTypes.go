package jobservertypes

/*
type NewJobQueue struct {
	NewJobs []NewJob
}

type Workers struct {
	Info []WorkerInfo
}*/

type NewJob struct {
	ID      int32
	Case    string
	JobType string
	File    string
}

type WorkerInfo struct {
	Name           string
	Host           string
	Responsibility string
}
