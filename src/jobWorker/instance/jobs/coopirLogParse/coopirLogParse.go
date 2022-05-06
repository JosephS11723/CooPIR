package coopirLogParse

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/JosephS11723/CooPIR/src/jobWorker/config"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/resultTypes"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/worker"
	"github.com/buger/jsonparser"
)

// CoopirLogParse attempt to extract information from a coopir log
func CoopirLogParse(job *dbtypes.Job, resultChan chan worker.ResultContainer, returnChan chan string) error {
	// get information
	caseUUID := job.CaseUUID
	fileUUID := job.Files[0]

	// open reader for
	r, err := os.Open(config.WorkDir + "/" + caseUUID + "/" + fileUUID)
	if err != nil {
		log.Println(err)
		return err
	}
	defer r.Close()

	// read bytes from file
	bytes, err := ioutil.ReadAll(r)

	// if there is an error
	if err != nil {
		// print the error
		log.Fatal(err)
	}

	var lastUploadedFileuuid string

	// for element in key "logs" using arrayEach
	jsonparser.ArrayEach(bytes, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		// unmarshal into map
		var m map[string]interface{}
		json.Unmarshal(value, &m)

		// create function for putting v into a pipe writer
		reader, writer := io.Pipe()

		// add to case
		go func() {
			// write v to pipe
			writer.Write(value)
			writer.Close()
		}()

		// convert time
		time := m["time"].(float64)

		// convert time to string
		timeString := strconv.Itoa(int(time))

		// create container
		resultChan <- worker.ResultContainer{
			JobResult: worker.JobResult{
				ResultType: resultTypes.CreateFile,
				JobUUID:    job.JobUUID,
				CaseUUID:   job.CaseUUID,
				Tags:       []string{"logEntry"},
				Relations:  []string{fileUUID + ":contains"},
				Name:       "logEntry" + ":" + timeString,
				Done:       false,
				FileUUID:   "none",
			},
			FileReader: reader,
		}
		//reader.Close()

		// get last uploaded uuid from return chan
		lastUploadedFileuuid = <-returnChan

		// for each key in map
		for k, v := range m {
			// create function for putting v into a pipe writer
			reader, writer := io.Pipe()

			// if key is "content", unmarshal into map
			if k == "content" {
				// if value is not nil
				if v != nil {
					// act based on type
					switch v.(type) {
					case string:
						// add to case
						go func() {
							// write v to pipe
							writer.Write([]byte(v.(string)))
							writer.Close()
						}()

						// create container
						resultChan <- worker.ResultContainer{
							JobResult: worker.JobResult{
								ResultType: resultTypes.CreateFile,
								JobUUID:    job.JobUUID,
								CaseUUID:   job.CaseUUID,
								Tags:       []string{k},
								Relations:  []string{lastUploadedFileuuid + ":contains"},
								Name:       k + ":" + v.(string),
								Done:       false,
								FileUUID:   "none",
							},
							FileReader: reader,
						}
						//reader.Close()

						// get last uploaded uuid from return chan
						<-returnChan
					case map[string]interface{}:
						// for key in content
						for kk, vv := range v.(map[string]interface{}) {
							// print key and value
							//log.Println(k, kk, vv)
							// add to case
							go func() {
								// write v to pipe
								writer.Write([]byte(vv.(string)))
								writer.Close()
							}()

							// create container
							resultChan <- worker.ResultContainer{
								JobResult: worker.JobResult{
									ResultType: resultTypes.CreateFile,
									JobUUID:    job.JobUUID,
									CaseUUID:   job.CaseUUID,
									Tags:       []string{kk, "logContent"},
									Relations:  []string{lastUploadedFileuuid + ":contains"},
									Name:       k + ":" + vv.(string),
									Done:       false,
									FileUUID:   "none",
								},
								FileReader: reader,
							}
							//reader.Close()

							// void last uploaded uuid from return chan
							<-returnChan
						}
					case int64:
						var vtype string

						// if v is an int64, convert to a string
						switch v.(type) {
						case int64:
							vtype = strconv.Itoa(int(v.(int64)))
						case float64:
							vtype = v.(string)
						}

						// add to case
						go func() {
							// write v to pipe
							writer.Write([]byte(vtype))
							writer.Close()
						}()

						// create container
						resultChan <- worker.ResultContainer{
							JobResult: worker.JobResult{
								ResultType: resultTypes.CreateFile,
								JobUUID:    job.JobUUID,
								CaseUUID:   job.CaseUUID,
								Tags:       []string{k},
								Relations:  []string{lastUploadedFileuuid + ":contains"},
								Name:       k + ":" + vtype,
								Done:       false,
								FileUUID:   "none",
							},
							FileReader: reader,
						}
						//reader.Close()

						// get last uploaded uuid from return chan
						<-returnChan
					}
				}
			} else {
				var vtype string

				// if v is an int64, convert to a string
				switch v.(type) {
				case int64:
					vtype = strconv.Itoa(int(v.(int64)))
				case float64:
					vtype = strconv.Itoa(int(v.(float64)))
				case string:
					vtype = v.(string)
				}

				// add to case
				go func() {
					// write v to pipe
					writer.Write([]byte(vtype))
					writer.Close()
				}()

				// create container
				resultChan <- worker.ResultContainer{
					JobResult: worker.JobResult{
						ResultType: resultTypes.CreateFile,
						JobUUID:    job.JobUUID,
						CaseUUID:   job.CaseUUID,
						Tags:       []string{k},
						Relations:  []string{lastUploadedFileuuid + ":contains"},
						Name:       k + ":" + vtype,
						Done:       false,
						FileUUID:   "none",
					},
					FileReader: reader,
				}
				//reader.Close()

				// get last uploaded uuid from return chan
				<-returnChan
			}
		}
	}, "logs")

	return nil
}
