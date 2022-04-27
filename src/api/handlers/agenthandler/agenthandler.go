package agenthandler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/JosephS11723/CooPIR/src/api/lib/coopirutil"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

var ChanMap = make(map[string]chan Work)

var Agents = make(map[string]*ClientInfo)

var AgentMutex sync.Mutex = sync.Mutex{}

// info the client sends us to identify itself
type ClientInfo struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

// work we are given from drew to send to the client. client responds with file and we put it in seaweed. we then add an entry in the database
type Work struct {
	Task string `json:"task"`
	// other file data here. this is what is put in the database
}

func AgentHandler(con *gin.Context) {
	// upgrade connection
	//c, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	c, err := upgrader.Upgrade(con.Writer, con.Request, nil)

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

	go func() {
		AgentMutex.Lock()
		defer AgentMutex.Unlock()
		Agents[clientInfo.UUID] = &clientInfo
	}()

	// create a channel for this client in the map
	ChanMap[clientInfo.UUID] = make(chan Work, 100)

	// defer closing the channel
	defer delete(ChanMap, clientInfo.UUID)
	defer delete(Agents, clientInfo.UUID)

	// TODO: THIS IS WHERE THE CONNECTION WOULD LOITER UNTIL WE HAVE WORK TO DO
	// TEST CODE: add a work item to the channel
	//ChanMap[clientInfo.UUID] <- Work{Task: "getlogs"}

	// channel for closing
	closeChan := make(chan bool, 1)

	var work Work

	// automatically handle ping pong
	go func() {
		for {
			// get file from client
			// get reader
			mt, reader, err := c.NextReader()

			if err != nil {
				log.Println("next reader:", err)
				log.Println(clientInfo.UUID, "Client disconnected")
				return
			}

			// switch on message type
			switch mt {
			case websocket.TextMessage:
				log.Println(clientInfo.UUID, "TextMessage")
				// read file and save (call function)
				readAndSaveFile(clientInfo.UUID, reader)

				log.Println("File written to disk")

				// print the work struct (DEBUG)
				log.Printf("Work: %+v", work)

				// THIS IS WHERE THE WORK WOULD BE ADDED TO THE DATABASE

				// send response to client
				err = c.WriteMessage(websocket.TextMessage, []byte("200"))

				if err != nil {
					log.Println(clientInfo.UUID, "send response write:", err)
					return
				}

				log.Println(clientInfo.UUID, "Response sent")
			case websocket.BinaryMessage:
				log.Println(clientInfo.UUID, "BinaryMessage")
				// read file and save (call function)
				readAndSaveFile(clientInfo.UUID, reader)

				log.Println("File written to disk")

				// send response to client
				err = c.WriteMessage(websocket.TextMessage, []byte("200"))

				if err != nil {
					log.Println(clientInfo.UUID, "send response write:", err)
					return
				}

				log.Println(clientInfo.UUID, "Response sent")
			case websocket.CloseMessage:
				log.Println(clientInfo.UUID, "CloseMessage")
				closeChan <- true
				return
			case websocket.PingMessage:
				log.Println(clientInfo.UUID, "PingMessage")
				err = c.WriteMessage(websocket.PongMessage, []byte("pong"))

				if err != nil {
					log.Println(clientInfo.UUID, "send pong write:", err)
					return
				}
			case websocket.PongMessage:
				log.Println(clientInfo.UUID, "PongMessage")
			}

		}
	}()

	// infinite get work and send it loop
	for {
		// wait for incoming work or get close from close channel (switch)
		select {
		case work = <-ChanMap[clientInfo.UUID]:
			log.Printf("Work: %+v", work)
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
		case <-closeChan:
			log.Printf("Client %s disconnected", clientInfo.UUID)
			return
		}

	}
}

func readAndSaveFile(fileName string, reader io.Reader) {
	// TODO: this is where we upload the reader to seaweed and get data. skip this step
	// TEST CODE: write to a local file
	file, err := os.Create(fileName)

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
}

func SubmitWork(c *gin.Context) {
	// read the the work from the params
	// get task from params
	params, _, err := coopirutil.ParseParams([]string{"task", "uuid"}, c.Request.URL.Query())

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ChanMap[params["uuid"]] <- Work{Task: params["task"]}
}

func GetAgents(c *gin.Context) {
	defer AgentMutex.Unlock()
	AgentMutex.Lock()
	c.JSON(http.StatusOK, Agents)
}
