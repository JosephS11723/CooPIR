package seaweed

import (
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/JosephS11723/CooPIR/src/jobWorker/config"
)

type SWMount struct {
	// Path is the path to the file on the seaweed server
	LocalPath  string `json:"localpath"`
	Remotepath string `json:"remotepath"`
}

// terminal command to mount
var mountCommand string = "mount -filer=" + config.FilerAddress + ":" + config.FilerPort + " -dir=\""

var unmountCommand string = "-f "

// CreateSWMount creates a seaweed mount given a file path
func CreateSWMount(localPath string, remotepath string) SWMount {
	return SWMount{
		LocalPath:  localPath,
		Remotepath: remotepath,
	}
}

// Mount initiates the mount
func (s *SWMount) Mount() error {
	// make mount string
	mountString := mountCommand + s.LocalPath + "\" -filer.path=\"" + s.Remotepath + "\""

	log.Println("Mounting: ", mountString)

	// mount
	return runCMD("./weed", strings.Split(mountString, " ")...)
}

// Unmount unmounts the mount
func (s *SWMount) Unmount() error {
	// make unmount string
	unmountString := unmountCommand + s.LocalPath

	// unmount
	return runCMD("unmount", strings.Split(unmountString, " ")...)
}

// runCMD runs a command in the terminal
func runCMD(cmd string, args ...string) error {
	execCmd := exec.Command(cmd, args...)
	return execCmd.Run()
}

// this creates a quick and dirty mount for all the files
func MountAllFiles() {
	for {
		err := runCMD("./weed", "mount", "-filer=filer:8888", "-filer.path=/files")
		if err != nil {
			log.Println("Error mounting: ", err)
		}

		// it must have disconnected, reconnect after 5 seconds
		time.Sleep(time.Duration(5) * time.Second)
	}
}
