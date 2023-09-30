package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
)

type JsonRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int           `json:"id"`
	Desc    string
}

func TestJsonRpcRequest(t *testing.T) {
	url := HTTPRPCLocalhost
	JsonRequestArr := make([]JsonRequest, 0, 0)

	// Echo
	// Test_Echo
	JsonRequestArr = append(JsonRequestArr, JsonRequest{
		Jsonrpc: "2.0",
		Method:  "test_Echo",
		Params: []interface{}{
			"f",
			3,
			map[string]string{"S": "v"},
		},
		Id:   1,
		Desc: "Test_echo",
	})

	fmt.Println("JSON START -----------------------")
	jsonBytes, _ := json.Marshal(JsonRequestArr)
	fmt.Println(string(jsonBytes))
	fmt.Println("JSON END -----------------------")

	//JsonRequestArr = append(JsonRequestArr, JsonRequest{
	//	Jsonrpc: "2.0",
	//	Method:  "test_peerInfo",
	//	Id:      1,
	//	Desc:    "PeerInfo",
	//})
	//JsonRequestArr = append(JsonRequestArr, JsonRequest{
	//	Jsonrpc: "2.0",
	//	Method:  "test_sleep",
	//	Params:  []inter{}{3},
	//	Id:      1,
	//	Desc:    "sleep",
	//})
	//he will block the goroutine
	//JsonRequestArr = append(JsonRequestArr, JsonRequest{
	//	Jsonrpc: "2.0",
	//	Method:  "test_block",
	//	Id:      1,
	//	Desc:    "Block",
	//})
	//JsonRequestArr = append(JsonRequestArr, JsonRequest{
	//	Jsonrpc: "2.0",
	//	Method:  "test_repeat",
	//	Params:  []inter{}{1, "x"},
	//	Id:      1,
	//	Desc:    "Repeat",
	//})

	for _, request := range JsonRequestArr {
		jsonData, err := json.Marshal(request)
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to send POST request: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Fatalf("request.Desc + \":\"+Failed to read response body: %v", err)
		} else {
			log.Println(request.Desc + ":" + string(body))
		}

	}

}

func TestSomeSubscription(t *testing.T) {
	u := url.URL{Scheme: "ws", Host: RPCLoclahost, Path: "/ws"}

	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	// Check response headers from the initial handshake
	fmt.Println("Response headers:")
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "nftest_subscribe",
		"params":  []interface{}{"someSubscription", 1, 2},
	}

	err = conn.WriteJSON(request)
	if err != nil {
		t.Fatalf("write: %v", err)
	}

	var response map[string]interface{}
	err = conn.ReadJSON(&response)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	subscriptionID, ok := response["result"].(string)
	if !ok {
		t.Fatalf("unexpected result: %v", response)
	}
	log.Println(subscriptionID)
	t.Skip()
}
