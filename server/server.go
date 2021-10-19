package server

import (
	"blockchain-prototype/core"
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Server struct {
	buf        *bytes.Buffer
	infoLogger *log.Logger
	Location   string
	state      *core.BlockChainState
}

func CreateServer(location string) Server {
	server := Server{
		state:    core.CreateBlockChainState(),
		Location: location,
	}
	server.buf = new(bytes.Buffer)
	server.infoLogger = log.New(server.buf, "BLOCKCHAIN: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	server.infoLogger.Printf("Stated the server")
	fmt.Print(server.buf)
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
