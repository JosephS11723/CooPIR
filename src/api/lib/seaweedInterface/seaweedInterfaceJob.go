package seaweedInterface

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/JosephS11723/CooPIR/src/api/config"
	libcrypto "github.com/JosephS11723/CooPIR/src/api/lib/crypto"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
)

// POSTFileJob allows the job result processing function to attempt to upload a job to the seaweed filer.
// This function automatically handles file conflicts and returns the UUID of the conflicting file in that event.
// If there is a conflict, it will be returned as an error. Handle it!
func POSTFileJob(caseUUID string, r io.Reader) (string, error) {
	// set filename to randomly generated name. change after hash operation
	// Use MakeUuid from dbInterface to ensure unique filename
	filename, err := dbInterface.MakeUuid()

	// error if failed to generate uuid
	if err != nil {
		return "", err
	}

	// create pipes
	md5Reader, md5Writer := io.Pipe()
	sha1Reader, sha1Writer := io.Pipe()
	sha256Reader, sha256Writer := io.Pipe()
	sha512Reader, sha512Writer := io.Pipe()
	POSTReader, POSTWriter := io.Pipe()

	// make channel for errors
	errChan := make(chan error, config.UpTaskCount)
	defer close(errChan)

	// make channel for hash result
	hashmd5Chan := make(chan []byte, 1)
	defer close(hashmd5Chan)

	// make channel for hash result
	hashsha1Chan := make(chan []byte, 1)
	defer close(hashsha1Chan)

	// make channel for hash result
	hashsha256Chan := make(chan []byte, 1)
	defer close(hashsha256Chan)

	// make channel for hash result
	hashsha512Chan := make(chan []byte, 1)
	defer close(hashsha512Chan)

	// sync group for synchronizing tasks
	var doWait sync.WaitGroup

	// addthe number of tasks to do
	doWait.Add(config.UpTaskCount)

	// spawn all the task goroutines. These look identical to
	// the TeeReader example, but pulled out into separate
	// methods for clarity
	go libcrypto.MD5FromReaderAsync(md5Reader, &doWait, errChan, hashmd5Chan)
	go libcrypto.Sha1FromReaderAsync(sha1Reader, &doWait, errChan, hashsha1Chan)
	go libcrypto.Sha256FromReaderAsync(sha256Reader, &doWait, errChan, hashsha256Chan)
	go libcrypto.Sha512FromReaderAsync(sha512Reader, &doWait, errChan, hashsha512Chan)
	go POSTFile(filename, caseUUID, POSTReader, &doWait, errChan)

	go func() {
		// after completing the copy, we need to close
		// the PipeWriters to propagate the EOF to all
		// PipeReaders to avoid deadlock
		defer md5Writer.Close()
		defer sha1Writer.Close()
		defer sha256Writer.Close()
		defer sha512Writer.Close()
		defer POSTWriter.Close()

		// build the multiwriter for all the pipes
		mw := io.MultiWriter(md5Writer, sha1Writer, sha256Writer, sha512Writer, POSTWriter)

		// copy the data into the multiwriter
		io.Copy(mw, r)
	}()

	// read hash output from functions
	filemd5Hash := <-hashmd5Chan
	filesha1Hash := <-hashsha1Chan
	filesha256Hash := <-hashsha256Chan
	filesha512Hash := <-hashsha512Chan

	// if logging for upload is enabled, print hashes to log
	if config.DoUploadLogging {
		fmt.Println("MD5:   ", hex.EncodeToString(filemd5Hash))
		fmt.Println("SHA1:  ", hex.EncodeToString(filesha1Hash))
		fmt.Println("SHA256:", hex.EncodeToString(filesha256Hash))
		fmt.Println("SHA512:", hex.EncodeToString(filesha512Hash))
	}

	// wait until all tasks are done
	doWait.Wait()

	// check errors
	for i := 0; i < config.UpTaskCount; i++ {
		hashErr := <-errChan
		if err != nil {
			// upload failed, panic!
			err = hashErr
		}
	}

	if err != nil {
		// upload failed
		return "", errors.New("upload to seaweed filer failed")
	}

	// check mongo for file existence and remove if duplicate
	conflictUUID, err := dbInterface.FindFileByHash(hex.EncodeToString(filesha512Hash), caseUUID)

	for err == nil {
		// file already exists, remove it
		err = DELETEFile(filename, caseUUID)
		if err != nil {
			// print error and retry
			log.Println("Error deleting file in seaweed:", err)
			continue
		}

		// file deleted, return conflict uuid
		return conflictUUID, nil
	}

	// return the uuid
	return filename, nil
}