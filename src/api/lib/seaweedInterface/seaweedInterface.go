package seaweedInterface

import (
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"sync"

	//"fmt"

	"github.com/JosephS11723/CooPIR/src/api/config"
	"github.com/gin-gonic/gin"
)

var filerAddress string = "http://filer"
var filerPort string = "8888"

// GETs a file from seaweed and streams it to the ctx writer
func GETFile(filename string, casename string, c *gin.Context) error {
	// create http client
	client := &http.Client{}

	// var for response
	var resp *http.Response

	// create request
	req, err := http.NewRequest(http.MethodGet, filerAddress+":"+filerPort+"/files/" + casename + "/" + filename, nil)
	if err != nil {
		return err
	}

	// do request
	resp, err = client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}

	// defer closing body
	defer resp.Body.Close()

	// set header
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
	c.Writer.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

	// copy file
	io.Copy(c.Writer, resp.Body)

	return nil
}

// PUTs a file to seaweed using the io reader passed
func POSTFile(filename string, casename string, r io.Reader, c *gin.Context, ss *sync.WaitGroup, errChan chan error) {
	// defer job finish
	defer ss.Done()

	var err error

	rr, w := io.Pipe()
	mpw := multipart.NewWriter(w)
	go func() {
		var part io.Writer
		defer w.Close()

		if part, err = mpw.CreateFormFile("file", filename); err != nil {
			log.Println(err)
			errChan <- err
			return
		}
		part = io.MultiWriter(part)
		if _, err = io.Copy(part, r); err != nil {
			log.Println(err)
			errChan <- err
			return
		}
		if err = mpw.Close(); err != nil {
			log.Println(err)
			errChan <- err
			return
		}
	}()

	var fileStat string

	// determine read-only status
	if config.ReadOnlyFiles {
		fileStat = "?mode=1555"
	} else {
		fileStat = "?mode=0777"
	}

	// create request
	resp, err := http.Post(filerAddress+":"+filerPort+"/files/" + casename + "/" + fileStat, mpw.FormDataContentType(), rr)
	log.Println(filerAddress + ":" + filerPort + "/files/" + casename + "/" + fileStat)
	if err != nil {
		log.Println(err)
		errChan <- err
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		errChan <- err
		return
	}
	//fmt.Print(string(ret))

	errChan <- nil
}

// DELETEs a file on seaweed given its name
func DELETEFile(filename string, c *gin.Context) error {
	// create http agent
	client := &http.Client{}

	// create request
	req, err := http.NewRequest(http.MethodDelete, filerAddress+":"+filerPort+"/files/"+filename, nil)

	if err != nil {
		log.Println(err)
		return err
	}

	// do request
	_, err = client.Do(req)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
