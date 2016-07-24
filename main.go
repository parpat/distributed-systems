package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

type LivePeers struct {
	Container_ID string
	Addr         string
}

func main() {

	//Initialize Server
	l, err := net.Listen("tcp", ":7575")
	fmt.Println("Listening")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	go listener(l)
	
	//Initialize client
	file, err := os.Open("livePeers.json")
	if err != nil {
		fmt.Println("error:", err)
	}
	decoder := json.NewDecoder(file)
	var peers []LivePeers
	err = decoder.Decode(&peers)
	if err != nil {
		fmt.Println("error:", err)
	}

	//time.Sleep(time.Second * 20)
	
	go func() {
		for{
		time.Sleep(time.Second * 23)
		rand.Seed(int64(time.Now().Second()))
		addr := peers[rand.Intn(len(peers))].Addr

		go client_routine(addr)
		fmt.Println("Started go routine for :", addr)
		}

	}()

	for {

		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Handle the connection in a new goroutine.
		go serveConn(conn)
	}

}

func listener(l net.Listener) {
	for {

		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Handle the connection in a new goroutine.
		go serveConn(conn)
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
	fmt.Println("Receiving data from: ", c.RemoteAddr)
	fmt.Println(resp)

	fmt.Fprintf(c, "Destination received data\n")

	time.Sleep(2 * time.Second)

}

func client_routine(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprint(conn, "Sending data from client!! ..\n")
	if err != nil {
		log.Print(err)
	}
	//fmt.Printf("bytes sent %d", n)

	//Receiving data
	for {
		resp, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Print(err)
			break
		}
		fmt.Println(resp)
	}

	conn.Close()
	fmt.Println("conn closed")
	//fmt.Println(d)
}
