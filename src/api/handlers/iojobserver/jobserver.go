package iojobserver

import (
	"github.com/gin-gonic/gin"
)

func CreateNewJob(queue_chan chan<- interface{}) gin.HandlerFunc {

	handler_func := func(c *gin.Context) {

	}

	return gin.HandlerFunc(handler_func)
}
