package worker

import (
	"net/url"
	"strconv"
	"strings"
)

// jobResultToParams converts a job result to a query encoded parameters string
func jobResultToParams(result *JobResult) string {
	// encode the struct as a map
	params := map[string]string{
		"jobuuid":    result.JobUUID,
		"name":       result.Name,
		"tags":       strings.Join(result.Tags, ","),
		"relations":  strings.Join(result.Relations, ","),
		"fileuuid":   result.FileUUID,
		"done":       strconv.FormatBool(result.Done),
		"resulttype": result.ResultType,
	}

	// convert the map to a query encoded string
	paramsStr := url.Values{}
	for k, v := range params {
		paramsStr.Set(k, v)
	}

	// set the params
	return paramsStr.Encode()
}

func jobTypesToParams(types []string) string {
	// encode the struct as a map
	params := map[string]string{
		"jobTypes": strings.Join(types, ","),
	}

	// convert the map to a query encoded string
	paramsStr := url.Values{}
	for k, v := range params {
		paramsStr.Set(k, v)
	}

	// set the params
	return paramsStr.Encode()
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