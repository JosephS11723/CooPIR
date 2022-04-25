package worker

// configuration for the paths of the api for each task

// basePath is the base path for the api
var basePath string = "/api/v1/jobs"

// path for obtaining work
var getWorkPath string = basePath + "/getwork"

// path for submiting results
var submitResultPath string = basePath + "/{jobuuid}/result"

// path for obtaining files related to a job
var getJobPath string = basePath + "/{jobuuid}/{fileuuid}"

// path for obtaining the status of a job
var getJobStatusPath string = basePath + "/status"

// path for registering the worker
var registerPath string = basePath + "/worker/new"