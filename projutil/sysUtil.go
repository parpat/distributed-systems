package projutil

import (
	"bytes"
	"log"
	"os/exec"
)

//GetHostInfo returns the ID and IP of the host container
//using the os command
func GetHostInfo() (string, string) {
	hostIP, err := exec.Command("hostname", "-i").Output()
	if err != nil {
		log.Fatal(err)
	}
	hostIP = bytes.TrimSuffix(hostIP, []byte("\n"))

	hostName, err := exec.Command("hostname").Output()
	if err != nil {
		log.Fatal(err)
	}
	hostName = bytes.TrimSuffix(hostName, []byte("\n"))

	return string(hostName), string(hostIP)
}
