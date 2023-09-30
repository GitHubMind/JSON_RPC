package main

import (
	"log"
)

const RPCLoclahost = "127.0.0.1:3000"
const HTTPRPCLocalhost = "http://" + RPCLoclahost + "/rpc"

// TODO depent config set
const GateWayLocalhost = "127.0.0.1:3000"

func setup() {
	log.Println("plz sure you already running RPC_SERVER and GateWay firstly!")
}
func teadown() {
	log.Println("test ends")

}

//func TestMain(m *testing.M) {
//	setup()
//	m.Run()
//	teadown()
//}
