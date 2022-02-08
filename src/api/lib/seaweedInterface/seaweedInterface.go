package seaweedInterface

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GETs a file from seaweed and streams it to the ctx writer
func GetFile(filename string, ctx gin.Context) error {
	// create http client
	client := &http.Client{}

	// var for response
	var resp *http.Response

	// create request
	req, err := http.NewRequest(http.MethodGet, "http://filer/" + filename, nil)
	if err != nil {
		return err
	}

	// do request
	resp, err = client.Do(req)

	if err != nil || response.StatusCode != http.StatusOK {
		return err
	}
	// defer closing body
	defer resp.Body.Close()

	ctx.Writer.Header().Set("Content-Disposition", "attachment; filename=" + filename)
	//ctx.Writer.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	io.Copy(ctx.Writer, resp.Body)

	return nil
}

// PUTs a file to seaweed using the io reader passed
func UploadFile(filename string, r io.Reader) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, "http://filer/" + filename, nil)
	if err != nil {
		log.Panicln(err)
	}

	_, err = client.Do(req)

	if err != nil {
		log.Panicln(err)
	}
}

// DELETEs a file on seaweed given its name
func DeleteFile(filename string) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, "http://filer/" + filename, nil)
	if err != nil {
		log.Panicln(err)
	}

	_, err = client.Do(req)

	if err != nil {
		log.Panicln(err)
	}
}