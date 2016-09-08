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
const TTL time.Duration = (time.Second) * 10

//REFRESHSEC to refresh /leader
const REFRESHSEC time.Duration = (time.Second) * 5

//API to interact with etcd
var kapi client.KeysAPI

//Leader of the cluster
var Leader string

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
		Leader = HostName
		go refreshLeader()
	}
}

//LeaderWatcher watches the /leader key for changes by blocking
func LeaderWatcher() {
	//initial check
	resp, err := kapi.Get(context.Background(), "/leader", nil)
	if err != nil {
		clierr := err.(client.Error)
		log.Println(clierr.Code)
		SetLeader()
	} else {
		Leader = resp.Node.Value
	}

	//keep watching for changes
	watcher := kapi.Watcher("/leader", nil)
	for {
		resp, err := watcher.Next(context.Background())
		if err != nil {
			log.Println(err.Error())
		}

		if resp.Action == "expire" {
			SetLeader()
		} else {
			Leader = resp.Node.Value
			log.Printf("Current Leader: %s\n", Leader)
		}
	}
}

func refreshLeader() {
	setopts := &client.SetOptions{PrevExist: "true", TTL: TTL, Refresh: true}
	for {
		//log.Println("Refreshing leader key")
		_, err := kapi.Set(context.Background(), "/leader", "", setopts)
		if err != nil {
			clierr := err.(client.Error)
			log.Printf("Leader refresh failed: %s\n", clierr.Message)
		}
		time.Sleep(REFRESHSEC)
	}
}

//SetPeerInfo sets this peer's current host(container)name and IP address
func SetPeerInfo(name, addr string) {
	_, err := kapi.Set(context.Background(), "/peers/"+name, addr, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("PeerInfo  registered to ETCD ")
	}

}

//GetPeers retrieves the peers list from etcd
func GetPeers() []Peer {
	var peers []Peer

	getopt := &client.GetOptions{Recursive: true, Sort: true, Quorum: true}
	resp, err := kapi.Get(context.Background(), "/peers", getopt)
	if err != nil {
		log.Println("Failed to obtain peer: ", err)
	} else {
		log.Println("Refreshed peer list")
		if (resp.Node).Nodes != nil {
			for _, node := range resp.Node.Nodes {
				peerName := strings.TrimPrefix(node.Key, "/peers/")
				peers = append(peers, Peer{Name: peerName, Addr: node.Value})
				log.Printf("Key: %q  Value: %q\n", peerName, node.Value)
			}

		}
	}

	return peers
}
