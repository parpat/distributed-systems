package main

import (
	"log"
	"strings"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

//ENDPOINT1 to etcd clusters
const ENDPOINT1 string = "http://172.17.0.1:2379"

//TTL for key
const TTL time.Duration = (time.Second) * 45

//API to interact with etcd
var kapi client.KeysAPI

func init() {
	cfg := client.Config{
		Endpoints: []string{ENDPOINT1},
		Transport: client.DefaultTransport,
		//Target endpoint timeout
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kapi = client.NewKeysAPI(c)
}

//SetLeader sets the leader key value
func SetLeader() {
	setopts := &client.SetOptions{PrevExist: "false", TTL: TTL}
	log.Println("Attempting to conquer!")
	_, err := kapi.Set(context.Background(), "/leader", HostName, setopts)
	if err != nil {
		clierr := err.(client.Error)
		log.Printf("Conquer failed: %s\n", clierr.Message)
	} else {
		log.Println("Conquer succeeded!")
		//go refreshLeader()
	}
}

//LeaderWatcher watches the /leader key for changes by blocking
func LeaderWatcher() {
	watcher := kapi.Watcher("/leader", nil)
	for {
		resp, err := watcher.Next(context.Background())
		if err != nil {
			log.Println(err.Error())
		}

		if resp.Action == "expire" {
			SetLeader()
		} else {
			log.Printf("Current Leader: %s\n", resp.Node.Value)
		}
	}
}

//SetPeerInfo sets this peer's current host(container)name and IP address
func SetPeerInfo(name, addr string) {
	resp, err := kapi.Set(context.Background(), "/peers/"+name, addr, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		// print common key info
		log.Printf("SetPeerInfo done. Metadata is %q\n", resp)
	}

}

//GetPeers retrieves the peers list from etcd
func GetPeers() []Peer {
	var peers []Peer

	getopt := &client.GetOptions{Recursive: true, Sort: true, Quorum: true}
	resp, err := kapi.Get(context.Background(), "/peers", getopt)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Response metadata: %q\n", resp)
		log.Printf("%q's Value: %q\n", resp.Node.Key, resp.Node.Value)

		if (resp.Node).Nodes != nil {
			for _, node := range resp.Node.Nodes {
				peerName := strings.TrimPrefix(node.Key, "/peers/")
				peers = append(peers, Peer{Name: peerName, Addr: node.Value})
				log.Printf("Key: %q  Value: %q", peerName, node.Value)
			}

		}
	}

	return peers
}