package main

import (
	"blockchain-prototype/core"
	"blockchain-prototype/server"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"net/rpc"
)

func main() {
	pvKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")

	if err != nil {
		log.Fatal("dialing:", err)
	} else {
		log.Print("connected")
	}
	// pvKey.Curve.m
	rpcblock := server.RPCBlock{MinedBlock: nil}
	err = client.Call("Server.MineBlock", (*core.Address)(&pvKey.PublicKey), &rpcblock)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	// ser.MineBlock()

	// for {
	// }
}
