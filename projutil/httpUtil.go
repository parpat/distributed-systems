package projutil

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

//ETCDAddr is the root address to the etcd server
const ETCDAddr string = "http://172.17.0.1:2379/v2/keys" // "http://localhost:2379/v2/keys"

//HTTPRequestHeaders adds default values to request header
func HTTPRequestHeaders(req *http.Request, encodedData string) {
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))
	req.Header.Add("X-Content-Length", strconv.Itoa(len(encodedData)))
}

//PutLeaderRequest creates a request with the object and time to live
func PutLeaderRequest(value, ttl string) *http.Request {
	putValues := url.Values{}
	putValues.Set("value", value)
	putValues.Add("ttl", ttl)
	encoded := putValues.Encode()

	req, err := http.NewRequest("PUT", ETCDAddr+"/leader", bytes.NewBufferString(encoded))
	if err != nil {
		log.Fatal("NewRequest:", err)
	}

	HTTPRequestHeaders(req, encoded)
	return req

}

//PutPeerRequest creates a request to add a peer to the peer list
func PutPeerRequest(peerName, peerIP string) *http.Request {
	putValues := url.Values{}
	putValues.Set("value", peerIP)
	encoded := putValues.Encode()

	req, err := http.NewRequest("PUT", ETCDAddr+"/peers/"+peerName, bytes.NewBufferString(encoded))
	if err != nil {
		log.Fatal("NewRequest:", err)
	}

	HTTPRequestHeaders(req, encoded)
	return req

}

//SendClientRequest creates a client and sends the provided Request
//The response body is read, closed, and returned
func SendClientRequest(req *http.Request) []byte {
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return resBody
}
