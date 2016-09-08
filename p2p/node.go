package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os/exec"
	"time"
)

//PORT is where the server listens
const PORT string = ":7575"

//Peer holds the name and network address of a node
type Peer struct {
	Name string
	Addr string
}

//peers holds the peer names and address of peers
//registered on etcd
var peers []Peer

//HostName is container ID
var HostName string

//HostIP is container IP
var HostIP string

func init() {
	HostName, HostIP = GetHostInfo()
}

func main() {
	//Register to ETCD
	SetPeerInfo(HostName, HostIP+PORT)

	go LeaderWatcher()

	go refreshPeers()

	//This routine will randomly select a peer to contact every 25 seconds
	go func() {
		for {
			time.Sleep(time.Second * 11)
			rand.Seed(int64(time.Now().Second()))
			if len(peers) > 0 {
				addr := peers[rand.Intn(len(peers))].Addr

				if addr != (HostIP + PORT) {
					go clientRoutine(addr)
					fmt.Println("Started go routine for :", addr)
				}
			}
		}
	}()

	//Initialize Server
	l, err := net.Listen("tcp", PORT)
	fmt.Println("Listening")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	//Keep serving connections as new routines
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Handle the connection in a new goroutine.
		go serveConn(conn)
	}

}

func refreshPeers() {
	for {
		peers = GetPeers()
		time.Sleep(time.Minute * 2)
	}
}

func serveConn(c net.Conn) {
	defer c.Close()
	//fmt.Println("Serving connection")

	resp, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		log.Print(err)
		//break
	}
	fmt.Println("Receiving data from: ", c.RemoteAddr().String())
	fmt.Println(resp)

	fmt.Fprintf(c, "Destination received data\n")

	time.Sleep(2 * time.Second)

}

//clientRoutine will execute a call to the given address and close the channel
//to indicate termination of the thread
func clientRoutine(addr string) {
	defer func() {
		str := recover()
		fmt.Println(str)
	}()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
	}

	_, err = fmt.Fprint(conn, "Sending data from: "+HostIP+" !!..\n")
	if err != nil {
		log.Println(err)
	}
	//fmt.Printf("bytes sent %d", n)

	//Receiving data
	for {
		resp, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println(resp)
	}

	fmt.Println("conn closed")
	//fmt.Println(d)
}

//GetHostInfo returns the ID and IP of the host container
//using the os command
func GetHostInfo() (string, string) {
	HostIP, err := exec.Command("hostname", "-i").Output()
	if err != nil {
		log.Fatal(err)
	}
	HostIP = bytes.TrimSuffix(HostIP, []byte("\n"))

	HostName, err := exec.Command("hostname").Output()
	if err != nil {
		log.Fatal(err)
	}
	HostName = bytes.TrimSuffix(HostName, []byte("\n"))

	return string(HostName), string(HostIP)
}
