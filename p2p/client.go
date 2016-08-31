package main

import (
	"log"
	"strings"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

//Peer holds the name and network address of a node
type Peer struct {
	Name string
	Addr string
}

//ENDPOINT1 to etcd clusters
const ENDPOINT1 string = "http://localhost:2379"

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
func SetLeader(name string) {
	setopts := &client.SetOptions{TTL: TTL}

	//set key Value
	resp, err := kapi.Set(context.Background(), "/leader"+name, "wannabe", setopts)
	if err != nil {
		log.Fatal(err)
	} else {
		// print common key info
		log.Printf("Set is done. Metadata is %q\n", resp)
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
