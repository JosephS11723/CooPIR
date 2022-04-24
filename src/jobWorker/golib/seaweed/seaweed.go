package seaweed

import (
	"os/exec"
	"strings"

	"github.com/JosephS11723/CooPIR/src/jobWorker/config"
)

type SWMount struct {
	// Path is the path to the file on the seaweed server
	LocalPath  string `json:"localpath"`
	Remotepath string `json:"remotepath"`
}

// terminal command to mount
var mountCommand string = "weed mount -filer=" + config.FilerAddress + ":" + config.FilerPort + " -dir="

var unmountCommand string = "unmount -f "

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
	mountString := mountCommand + s.LocalPath + " -filer.path=" + s.Remotepath

	// mount
	return runCMD("sudo", strings.Split(mountString, " ")...)
}

// Unmount unmounts the mount
func (s *SWMount) Unmount() error {
	// make unmount string
	unmountString := unmountCommand + s.LocalPath

	// unmount
	return runCMD("sudo", strings.Split(unmountString, " ")...)
}

// runCMD runs a command in the terminal
func runCMD(cmd string, args ...string) error {
	execCmd := exec.Command(cmd, args...)
	return execCmd.Run()
}
