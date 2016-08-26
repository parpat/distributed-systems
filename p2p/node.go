package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/parth/projutil"
	"github.com/parth/projutil/parser"
)

//peerMap holds the peer names and address of peers
//registered on etcd
var peers []parser.Peer
var hostName, hostIP string

func main() {

	hostName, hostIP = projutil.GetHostInfo()
	req := projutil.PutPeerRequest(hostName, hostIP+":7575")
	response := projutil.SendClientRequest(req)
	fmt.Println("ETCD Registration: ", string(response))

	go refreshPeers()

	//This routine will randomly select a peer to contact every 25 seconds
	go func() {
		for {
			time.Sleep(time.Second * 11)
			rand.Seed(int64(time.Now().Second()))
			if len(peers) > 0 {
				addr := peers[rand.Intn(len(peers))].Addr

				go clientRoutine(addr)
				fmt.Println("Started go routine for :", addr)
			}
		}
	}()

	//Initialize Server
	l, err := net.Listen("tcp", ":7575")
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
		peers = getPeers()
		time.Sleep(time.Minute * 2)
	}
}

func getPeers() []parser.Peer {
	res, err := http.Get(projutil.ETCDAddr + "/peers/?recursive=true")
	if err != nil {
		log.Fatal("HTTP Get Peers: ", err)
	}
	resBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	peers := parser.PeerList(resBody)
	return peers
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

	_, err = fmt.Fprint(conn, "Sending data from: "+hostIP+" !!..\n")
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
