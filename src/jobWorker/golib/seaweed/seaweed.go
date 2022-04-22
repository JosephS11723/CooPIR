package seaweed

type SWMount struct {
	// Path is the path to the file on the seaweed server
	LocalPath  string `json:"localpath"`
	Remotepath string `json:"remotepath"`
}

// terminal command to unmount
var unmountCommand string = "sudo umount -f "

// terminal command to mount
var mountCommand string = "weed mount -filer=filer:8888 -dir=/some/existing/dir -filer.path=/one/remote/folder"

// CreateSWMount creates a seaweed mount given a file path
func CreateSWMount(localPath string, remotepath string) SWMount {
	return SWMount{
		LocalPath:  localPath,
		Remotepath: remotepath,
	}
}

// Mount initiates the mount
func (s *SWMount) Mount() error {
	// TODO: mount
	return nil
}

// Unmount unmounts the mount
func (s *SWMount) Unmount() error {
	// TODO: unmount
	return nil
}