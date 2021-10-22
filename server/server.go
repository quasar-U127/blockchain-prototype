package server

import (
	"blockchain-prototype/core/node"
	"blockchain-prototype/core/structure/block"
	"blockchain-prototype/core/structure/transaction"
	"blockchain-prototype/core/wallet"
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Nothing struct{}

type Server struct {
	buf        *bytes.Buffer
	infoLogger *log.Logger
	Location   string
	serverNode node.Node
	wal        wallet.Wallet
}

func CreateServer(location string) Server {
	// TODO: Initialize the server correctly
	server := Server{
		Location:   location,
		wal:        wallet.CreateWallet(),
		serverNode: node.CreateNode(),
	}

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

type RPCAddress struct {
	Add wallet.Address
}

func (ser *Server) GetAddress(name string, rpcAddress *RPCAddress) error {
	rpcAddress.Add = ser.wal.Get(name)
	return nil
}

type RPCBlock struct {
	MinedBlock *block.Block
}

func (ser *Server) MineBlock(name string, rpcBlock *RPCBlock) error {
	address := ser.wal.Get(name)
	status := ser.serverNode.Mine(&address)

	ser.buf.Reset()
	if status {
		rpcBlock.MinedBlock = ser.serverNode.GetBlock(ser.serverNode.GetEnd())
		ser.infoLogger.Printf("Mined Block %d : %x", rpcBlock.MinedBlock.Header.Height, rpcBlock.MinedBlock.Id())
		ser.infoLogger.Printf("With %d transactions", len(rpcBlock.MinedBlock.Txns))
	} else {
		ser.infoLogger.Print("Mining Failed")

	}
	fmt.Print(ser.buf)
	return nil
}

func (ser *Server) SubmitTransaction(txn *transaction.Transaction, status *bool) error {
	*status = ser.serverNode.AddTransaction(txn)
	return nil
}

func (ser *Server) CreateAddress(name string, status *bool) error {
	ser.wal.Add(name)
	*status = true
	return nil
}

type RPCAddressList struct{ Addresses map[string]string }

func (ser *Server) GetAddressList(nothing Nothing, list *RPCAddressList) error {
	for _, name := range ser.wal.GetAddressNames() {
		add := ser.wal.Get(name)
		list.Addresses[name] = add.ShortString(10)
	}
	return nil
}

type RPCTXOSet struct {
	Set map[transaction.OutPoint]transaction.Output
}

func (ser *Server) GetUTXO(add *wallet.Address, utxoSet *RPCTXOSet) error {
	for op, o := range ser.serverNode.GetUTXO(add) {
		utxoSet.Set[op] = o
	}
	return nil
}

func (ser *Server) PrintTXO(outpoint *transaction.OutPoint, txo *transaction.Output) string {
	return fmt.Sprintf("(%x,%d)->(%s,%d)", outpoint.Id[:2], outpoint.N, ser.wal.GetName(txo.Reciever), txo.Value)
}
