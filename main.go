package main

import (
	"blockchain-prototype/core"
	"blockchain-prototype/core/structure/transaction"
	"blockchain-prototype/server"
	"fmt"
)

func Menu(options []string) string {
	for i, option := range options {
		fmt.Printf("%d: %s\n", i, option)
	}
	var option int
	fmt.Scan(&option)
	return options[option]
}

func main() {
	serv := server.CreateServer("127.0.0.1:1234")
	ser := &serv
	for {
		fmt.Println("---Main Menu---")
		options := []string{"block", "transaction", "address", "exit"}

		option := Menu(options)
		switch option {

		case "block":
			{
				BlockMenu(ser)
			}
		case "transaction":
			{
				TransactionMenu(ser)
			}
		case "address":
			{
				AddressMenu(ser)
			}
		case "exit":
			{
				return
			}

		}

	}
	// server.Serve(&ser)
}

func BlockMenu(ser *server.Server) {
	for {
		fmt.Println("---Block Menu---")
		options := []string{"mine", "get", "return"}
		option := Menu(options)
		switch option {
		case "mine":
			{
				addList := server.RPCAddressList{Addresses: map[string]string{}}
				ser.GetAddressList(server.Nothing{}, &addList)

				fmt.Printf("Enter Address from ( ")
				for name := range addList.Addresses {
					fmt.Printf("%s, ", name)
				}
				fmt.Printf("):")
				var add string
				fmt.Scan(&add)
				if _, ok := addList.Addresses[add]; !ok {
					continue
				}
				block := server.RPCBlock{MinedBlock: nil}
				ser.MineBlock(add, &block)
				if block.MinedBlock == nil {
					continue
				}
				for _, txn := range block.MinedBlock.Txns {
					txn.Print()
				}
			}
		case "return":
			{
				return
			}
		}
	}
}

func AddressMenu(ser *server.Server) {
	for {
		fmt.Println("---Address Menu---")
		options := []string{"list", "new", "return"}
		option := Menu(options)
		switch option {
		case "list":
			{
				addList := server.RPCAddressList{Addresses: map[string]string{}}
				ser.GetAddressList(server.Nothing{}, &addList)
				for name, pbkey := range addList.Addresses {
					fmt.Printf("\t\t %s: %x\n", name, pbkey)
				}
			}
		case "new":
			{
				addList := server.RPCAddressList{Addresses: map[string]string{}}
				ser.GetAddressList(server.Nothing{}, &addList)

				fmt.Printf("Enter new address name except ( ")
				for name := range addList.Addresses {
					fmt.Printf("%s, ", name)
				}
				fmt.Printf("):")
				var add string
				fmt.Scan(&add)
				status := false
				ser.CreateAddress(add, &status)
			}
		case "return":
			{
				return
			}
		}
	}
}

func TransactionMenu(ser *server.Server) {
	for {
		fmt.Println("---Transaction Menu---")
		options := []string{"new", "return"}
		option := Menu(options)
		switch option {
		case "new":
			{
				CreateTransaction(ser)
			}
		case "return":
			{
				return
			}
		}
	}
}

func CreateTransaction(ser *server.Server) {
	set := server.RPCTXOSet{Set: map[transaction.OutPoint]transaction.Output{}}
	ser.GetUTXO(nil, &set)
	txoList := make([]transaction.OutPoint, 0, len(set.Set))
	for op := range set.Set {
		txoList = append(txoList, op)
	}

	fmt.Println()
	fmt.Println()
	for i, out := range txoList {
		txo := set.Set[out]
		fmt.Printf("%d : %s\n", i, ser.PrintTXO(&out, &txo))
	}
	fmt.Print("Number of inputs: ")
	inSize := 0
	fmt.Scan(&inSize)
	totalValue := 0
	fmt.Print("Enter Inputs\n")
	inputs := []transaction.OutPoint{}

	for i := 0; i < inSize; i++ {
		input := 0
		fmt.Scanf("%d", &input)
		inputs = append(inputs, txoList[input])
		totalValue += int(set.Set[txoList[input]].Value)
	}
	fmt.Printf("Total Amount : %d\n", totalValue)

	addList := server.RPCAddressList{Addresses: map[string]string{}}
	ser.GetAddressList(server.Nothing{}, &addList)

	fmt.Print("Number of Outputs: ")
	outSize := 0
	fmt.Scan(&outSize)
	totalOutValue := 0
	fmt.Print("Enter Outputs\n")
	outputs := []transaction.Output{}
	for i := 0; i < outSize; i++ {
		name := ""
		value := 0
		fmt.Scanf("%s %d", &name, &value)
		address := server.RPCAddress{}
		ser.GetAddress(name, &address)
		outputs = append(outputs, transaction.Output{Reciever: address.Add, Value: core.Cookies(value)})
		totalOutValue += value
	}
	status := true
	ser.SubmitTransaction(&transaction.Transaction{Inputs: inputs, Outputs: outputs}, &status)
}
