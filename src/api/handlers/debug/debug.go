package debug

import "github.com/gin-gonic/gin"

// Response to ping-pong challenge
func PingPong(c *gin.Context) {
	c.JSON(200, gin.H{
		"data": "pong",
	})
}
