package jobservertypes

type NewJob struct {
	request_id int32
	case_name  string
	job_type   string
	file       string
}

type WorkerInfo struct {
	Name           string
	Host           string
	Responsibility string
}
