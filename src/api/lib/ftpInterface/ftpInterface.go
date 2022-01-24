package ftpinterface

import (
	"io"
	"log"
	_ "sync"
	"time"

	"github.com/JosephS11723/CooPIR/src/api/config"
	"github.com/secsy/goftp"
)

//var FTPLock sync.RWMutex

// Connects to the ftp server
func FtpConnect() *goftp.Client {
	// connection configuration information
	ftpConfig := goftp.Config{
		User:               config.FtpUsername,
		Password:           config.FtpPassword,
		ConnectionsPerHost: config.FtpConnectionsPerHost,
		Timeout:            time.Duration(config.FtpTimeout) * time.Second,
	}

	// dial the server and login
	client, err := goftp.DialConfig(ftpConfig, config.FtpIP)
	if err != nil {
		log.Panicln(err)
	}

	// return the client pointer
	return client
}

// Closes the ftp server connection
func FtpClose(c *goftp.Client) {
	if err := c.Close(); err != nil {
		log.Panicln(err)
	}
}

// Retrieves a file from the ftp server and streams it into the io.writer
func ReadFile(c *goftp.Client, filename string, w io.Writer) {
	// lock read lock
	//FTPLock.RLock()

	// unlock read lock
	//defer FTPLock.RUnlock()

	// attempt file connection retreive
	err := c.Retrieve(filename, w)
	if err != nil {
		log.Panicln(err)
	}
}

// Streams (sends) a file from the io.writer to the ftp server
func WriteFile(c *goftp.Client, filename string, r io.Reader) error {
	// lock write lock
	//FTPLock.Lock()

	// unlock write lock
	//defer FTPLock.Unlock()

	// attempt file store
	err := c.Store(filename, r)
	if err != nil {
		log.Panicln(err)
		return err
	}

	// no error, return nil
	return nil
}

// Deletes a file from the ftp server based on the filename
func DeleteFile(c *goftp.Client, filename string) {
	// lock write lock
	//FTPLock.Lock()

	// unlock write lock
	//defer FTPLock.Unlock()

	// delete file
	err := c.Delete(filename)
	if err != nil {
		log.Panicln(err)
	}

}
