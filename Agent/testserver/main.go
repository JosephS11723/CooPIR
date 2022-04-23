package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:4201", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var ChanMap = make(map[string]chan Work)

// info the client sends us to identify itself
type ClientInfo struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

// work we are given from drew to send to the client. client responds with file and we put it in seaweed. we then add an entry in the database
type Work struct {
	Task     string `json:"task"`
	FileName string `json:"fileName"`
	// other file data here. this is what is put in the database
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	// TODO: add authentication here

	// print client connected
	log.Printf("client connected: %s", c.RemoteAddr())

	// Get message from client announcing client information
	_, message, err := c.ReadMessage()

	if err != nil {
		log.Println("read:", err)
		return
	}

	log.Printf("recv: %s", message)

	// unmarshal client information
	var clientInfo ClientInfo

	err = json.Unmarshal(message, &clientInfo)

	if err != nil {
		log.Println("unmarshal:", err)
		return
	}

	log.Printf("Client Info: %+v", clientInfo)

	// create a channel for this client in the map
	ChanMap[clientInfo.UUID] = make(chan Work, 10)

	// defer closing the channel
	defer delete(ChanMap, clientInfo.UUID)

	// THIS IS WHERE THE CONNECTION WOULD LOITER UNTIL WE HAVE WORK TO DO
	// TEST CODE: add a work item to the channel
	ChanMap[clientInfo.UUID] <- Work{Task: "GetLogs"}

	// infinite loop
	for {
		// wait for incoming work
		work := <-ChanMap[clientInfo.UUID]

		// marshal work
		workJSON, err := json.Marshal(work)

		if err != nil {
			log.Println("work marshal:", err)
			return
		}

		// send work to client
		err = c.WriteMessage(websocket.TextMessage, workJSON)

		if err != nil {
			log.Println("send work write:", err)
			return
		}

		log.Printf("sent: %s", workJSON)

		// get file from client
		// get reader
		_, reader, err := c.NextReader()

		if err != nil {
			log.Println("next reader:", err)
			return
		}

		// TODO: this is where we upload the reader to seaweed and get data. skip this step
		// TEST CODE: write to a local file
		file, err := os.Create("test.txt")

		if err != nil {
			log.Println("create file:", err)
			return
		}

		defer file.Close()

		_, err = io.Copy(file, reader)

		if err != nil {
			log.Println("copy:", err)
			return
		}

		log.Println("File written to disk")

		// send response to client
		err = c.WriteMessage(websocket.TextMessage, []byte("200"))

		if err != nil {
			log.Println("send response write:", err)
			return
		}

		log.Println("Response sent")
	}
}

func main() {
	// create memory for map
	ChanMap = make(map[string]chan Work)

	// create simple handler for websocket
	http.HandleFunc("/agent", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
