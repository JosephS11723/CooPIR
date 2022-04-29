package config

// seaweed config
// the ip address or domain of the seaweed filer
var FilerAddress string = "filer"

// the port of the seaweed filer
var FilerPort string = "8888"

// the name of the api container
var ApiName string = "api-server"

// the port the api listens on
var ApiPort string = "8080"

// GetJobInterval is the interval between getting jobs
var GetJobInterval int = 10

// SubmitResultInterval is the interval between submitting results
var SubmitResultInterval int = 0

// WorkDir is the directory where the files are stored
var WorkDir string = "/usr/src/app"