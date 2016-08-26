//Package parser parses json responses from etcd
package parser

import (
	"encoding/json"
	"fmt"
	"strings"
)

//Peer holds the name and network address of a node
type Peer struct {
	Name string
	Addr string
}

//PeerList Parses the node list portion of the List dir etcd json response.
func PeerList(jsonRaw []byte) []Peer {

	//etcdlist := []byte(`{"action":"get","node":{"key":"/peers","dir":true,"nodes":[{"key":"/peers/dumbpeerone","value":"152.145.25.26","modifiedIndex":55,"createdIndex":55},{"key":"/peers/dumbpeertwo","value":"152.145.25.27","modifiedIndex":56,"createdIndex":56}],"modifiedIndex":33,"createdIndex":33}}`)
	var jsonData interface{}

	err := json.Unmarshal(jsonRaw, &jsonData)
	if err != nil {
		fmt.Println(err)
	}

	parsedList := jsonData.(map[string]interface{})

	//Creating map to parse
	peerMap := make(map[string]string)

	//Parse to nested map, and find nodes
	node := parsedList["node"].(map[string]interface{})
	nodes := node["nodes"].([]interface{})
	//fmt.Printf("nodes is %T\n", node["nodes"])

	for _, item := range nodes {
		//Parsing to name, address peermap
		peer := item.(map[string]interface{})
		peerName := peer["key"].(string)

		//Trimming dir prefix "/peers"
		peerName = strings.TrimPrefix(peerName, "/peers/")
		peerMap[peerName] = peer["value"].(string)
	}

	//Converting to slice of peers
	var peers []Peer
	for name, addr := range peerMap {
		peers = append(peers, Peer{Name: name, Addr: addr})
	}

	// for name, addr := range PeerMap {
	// 	fmt.Printf("%[1]s: %[2]s\n", name, addr)
	// }

	return peers
}
