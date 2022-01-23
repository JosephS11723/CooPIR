package iodb

import (
	"log"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/gin-gonic/gin"
)

func DbTest(c *gin.Context) {
	log.Println("[DEBUG] iodb.DbTest()")
	// Get Client, Context, CalcelFunc and
	// err from connect method.
	client, ctx, cancel, err := dbInterface.DbConnect()
	if err != nil {
		log.Panicln(err)
	}

	// Ping mongoDB with Ping method
	dbInterface.DbPing(client, ctx)

	// Release resource when the main
	// function is returned.
	defer dbInterface.DbClose(client, ctx, cancel)
}
