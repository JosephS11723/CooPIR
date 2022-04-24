package httputil

import (
	"fmt"
	"net/url"
)

//parses the vals from c.Requests.URL.Query(); returns a map of keys to vals and/or a map of keys to slices of vals
func ParseParams(keys []string, vals url.Values) (map[string]string, map[string][]string, error) {

	var single_vals map[string]string
	var multi_vals map[string][]string

	for _, key := range keys {

		//if its just 1 val (like uuid = "1234"), then put into single val
		if len(vals[key]) == 1 {

			single_vals[key] = vals[key][0]

			//if its an array of vals (like files = ["1234", "5678"]), multi vals
		} else if len(vals[key]) > 1 {

			multi_vals[key] = vals[key]

			//missing, so just return an error
		} else {

			return nil, nil, fmt.Errorf("key %s not found", key)

		}

	}

	return single_vals, multi_vals, nil
}

//same as ParseParams, except it doesn't error if the val isn't there
func TryParseParams(keys []string, vals url.Values) (map[string]string, map[string][]string) {

	var single_vals map[string]string
	var multi_vals map[string][]string

	for _, key := range keys {

		//if its just 1 val (like uuid = "1234"), then put into single val
		if len(vals[key]) == 1 {

			single_vals[key] = vals[key][0]

			//if its an array of vals (like files = ["1234", "5678"]), multi vals
		} else if len(vals[key]) > 1 {

			multi_vals[key] = vals[key]

			//missing, so just return an error
		} else {

		}

	}

	return single_vals, multi_vals
}
