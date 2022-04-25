package worker

import (
	"net/url"
	"strconv"
)

// jobResultToParams converts a job result to a query encoded parameters string
func jobResultToParams(result *JobResult) string {
	paramsStr := Params{Values: url.Values{}}

	paramsStr.Add("jobuuid", result.JobUUID)
	paramsStr.Add("name", result.Name)
	paramsStr.AddArray("tags", result.Tags)
	paramsStr.AddArray("relations", result.Relations)
	paramsStr.Add("fileuuid", result.FileUUID)
	paramsStr.Add("done", strconv.FormatBool(result.Done))
	paramsStr.Add("resulttype", result.ResultType)
	paramsStr.Add("caseuuid", result.CaseUUID)

	// set the params
	return paramsStr.Values.Encode()
}

func jobTypesToParams(types []string) string {
	paramsStr := Params{Values: url.Values{}}

	paramsStr.AddArray("jobTypes", types)

	// set the params
	return paramsStr.Values.Encode()
}

func uuidToParams(uuid string) string {
	// encode the data as a map
	params := map[string]string{
		"uuid": uuid,
	}

	// convert the map to a query encoded string
	paramsStr := url.Values{}
	for k, v := range params {
		paramsStr.Set(k, v)
	}

	// set the params
	return paramsStr.Encode()
}

func registrationToParams(jobtype string, name string) string {
	paramsStr := Params{Values: url.Values{}}

	paramsStr.Add("jobtype", jobtype)
	paramsStr.Add("name", name)

	// set the params
	return paramsStr.Values.Encode()
}

type Params struct {
	Values url.Values
}

func (u *Params) Set(key, value string) {
	u.Values.Set(key, value)
}

func (u *Params) Add(key, value string) {
	u.Values.Add(key, value)
}

func (u *Params) AddArray(key string, values []string) {
	for _, value := range values {
		u.Add(key, value)
	}
}
