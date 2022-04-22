package seaweed

import (
	"github.com/JosephS11723/CooPIR/src/jobWorker/config"
)

type SWMount struct {
	// Path is the path to the file on the seaweed server
	LocalPath  string `json:"localpath"`
	Remotepath string `json:"remotepath"`
}

// terminal command to unmount
var unmountCommand string = "sudo umount -f "

// terminal command to mount
var mountCommand string = "weed mount -filer=" + config.FilerAddress + ":" + config.FilerPort + " -dir="

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
	// TODO: mount
	return nil
}

// Unmount unmounts the mount
func (s *SWMount) Unmount() error {
	// TODO: unmount
	return nil
}
