package iojobserver

import (
	"log"

	"github.com/JosephS11723/CooPIR/src/api/lib/jobservertypes"
	"github.com/gin-gonic/gin"
)

/*this function returns the handler that gin uses;
inside the handler function, the json body of the http request is bound to a NewJob
that gets pushed into the queue channel*/
func CreateNewJob(queue_chan chan<- jobservertypes.NewJob) gin.HandlerFunc {

	handler_func := func(c *gin.Context) {

		var job_request jobservertypes.NewJob

		err := c.BindJSON(job_request)

		if err != nil {
			log.Panicln(err)
		}

		queue_chan <- job_request
	}

	return gin.HandlerFunc(handler_func)
}
