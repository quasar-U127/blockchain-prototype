package server

import (
	"blockchain-prototype/core"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Nothing struct{}

type Server struct {
	buf          *bytes.Buffer
	infoLogger   *log.Logger
	Location     string
	state        *core.BlockChainState
	addresses    map[string]*ecdsa.PrivateKey
	addresToName map[ecdsa.PublicKey]string
}

func CreateServer(location string) Server {
	server := Server{
		state:    core.CreateBlockChainState(),
		Location: location,
	}
	server.addresses = map[string]*ecdsa.PrivateKey{}
	server.addresToName = map[ecdsa.PublicKey]string{}
	server.buf = new(bytes.Buffer)
	server.infoLogger = log.New(server.buf, "BLOCKCHAIN: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	server.infoLogger.Printf("Stated the server")
	fmt.Print(server.buf)
	status := true
	server.CreateAddress("root", &status)
	return server
}

func Serve(ser *Server) {
	rpc.Register(ser)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ser.Location)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

type RPCBlock struct {
	MinedBlock *core.Block
}

func (server *Server) MineBlock(address *core.Address, rpcBlock *RPCBlock) error {
	rpcBlock.MinedBlock = server.state.MineBlock(address)
	server.buf.Reset()
	if rpcBlock.MinedBlock != nil {
		server.infoLogger.Printf("Mined Block %d : %x", rpcBlock.MinedBlock.Header.Height, rpcBlock.MinedBlock.ComputeHash())
		server.infoLogger.Printf("With %d transactions", len(rpcBlock.MinedBlock.Txns))
	} else {
		server.infoLogger.Print("Mining Failed")

	}
	fmt.Print(server.buf)
	return nil
}

func (server *Server) SubmitTransaction(txn *core.Transaction, status *bool) error {
	*status = server.state.Frontier.VerifyTransaction(txn)
	if *status {
		server.state.InsertTransaction(txn)
	}
	return nil
}

func (server *Server) CreateAddress(name string, status *bool) error {
	*status = true
	var err error
	server.addresses[name], err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	server.addresToName[server.addresses[name].PublicKey] = name
	return err
}

type RPCAddressList struct{ Addresses map[string]*ecdsa.PublicKey }

func (server *Server) GetAddressList(nothing Nothing, list *RPCAddressList) error {
	for name, key := range server.addresses {
		list.Addresses[name] = &key.PublicKey
	}
	return nil
}
func (server *Server) GetUTXO(address *core.Address, utxoSet *core.TXOSet) error {
	if address == nil {
		*utxoSet = server.state.Frontier.Frontier
	} else {
		set := core.TXOSet{Set: map[core.OutPoint]core.TXO{}}
		for outpoint, txout := range server.state.Frontier.Frontier.Set {
			if txout.Reciever == *address {
				set.Set[outpoint] = txout
			}
		}
		*utxoSet = set
	}
	return nil
}

func (server *Server) PrintTXO(outpoint *core.OutPoint, txo *core.TXO) string {
	return fmt.Sprintf("(%x,%d)->(%s,%d)", outpoint.Id[:2], outpoint.N, server.addresToName[ecdsa.PublicKey(txo.Reciever)], txo.Value)
}
