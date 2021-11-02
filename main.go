package main

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	shell "github.com/ipfs/go-ipfs-api"

	u "github.com/ipfs/go-ipfs-util"
)

var sh *shell.Shell
var ncalls int

var _ = time.ANSIC

func sleep() {
	ncalls++
	//time.Sleep(time.Millisecond * 5)
}

func randString() string {
	alpha := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	l := rand.Intn(10) + 2

	var s string
	for i := 0; i < l; i++ {
		s += string([]byte{alpha[rand.Intn(len(alpha))]})
	}
	return s
}

func makeRandomObject() (string, error) {
	// do some math to make a size
	x := rand.Intn(120) + 1
	y := rand.Intn(120) + 1
	z := rand.Intn(120) + 1
	size := x * y * z

	r := io.LimitReader(u.NewTimeSeededRand(), int64(size))
	sleep()
	return sh.Add(r)
}

func makeRandomDir(depth int) (string, error) {
	if depth <= 0 {
		return makeRandomObject()
	}
	sleep()
	empty, err := sh.NewObject("unixfs-dir")
	if err != nil {
		return "", err
	}

	curdir := empty
	for i := 0; i < rand.Intn(8)+2; i++ {
		var obj string
		if rand.Intn(2) == 1 {
			obj, err = makeRandomObject()
			if err != nil {
				return "", err
			}
		} else {
			obj, err = makeRandomDir(depth - 1)
			if err != nil {
				return "", err
			}
		}

		name := randString()
		sleep()
		nobj, err := sh.PatchLink(curdir, name, obj, true)
		if err != nil {
			return "", err
		}
		curdir = nobj
	}

	return curdir, nil
}

func main() {
	sh = shell.NewShell("localhost:5001")
	msg := "final msg"
	var err error
	for i := 0; i < 12; i++ {
		msg, err = sh.Add(strings.NewReader(msg))
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(msg)
	}
	// for {
	// 	time.Sleep(time.Second * 1000)
	// }
}

// package main

// import (
// 	"blockchain-prototype/core"
// 	"blockchain-prototype/core/structure/transaction"
// 	"blockchain-prototype/server"
// 	"fmt"
// )

// func Menu(options []string) string {
// 	for i, option := range options {
// 		fmt.Printf("%d: %s\n", i, option)
// 	}
// 	var option int
// 	fmt.Scan(&option)
// 	return options[option]
// }

// func main() {
// 	serv := server.CreateServer("127.0.0.1:1234")
// 	ser := &serv
// 	for {
// 		fmt.Println("---Main Menu---")
// 		options := []string{"block", "transaction", "address", "exit"}

// 		option := Menu(options)
// 		switch option {

// 		case "block":
// 			{
// 				BlockMenu(ser)
// 			}
// 		case "transaction":
// 			{
// 				TransactionMenu(ser)
// 			}
// 		case "address":
// 			{
// 				AddressMenu(ser)
// 			}
// 		case "exit":
// 			{
// 				return
// 			}

// 		}

// 	}
// 	// server.Serve(&ser)
// }

// func BlockMenu(ser *server.Server) {
// 	for {
// 		fmt.Println("---Block Menu---")
// 		options := []string{"mine", "get", "return"}
// 		option := Menu(options)
// 		switch option {
// 		case "mine":
// 			{
// 				addList := server.RPCAddressList{Addresses: map[string]string{}}
// 				ser.GetAddressList(server.Nothing{}, &addList)

// 				fmt.Printf("Enter Address from ( ")
// 				for name := range addList.Addresses {
// 					fmt.Printf("%s, ", name)
// 				}
// 				fmt.Printf("):")
// 				var add string
// 				fmt.Scan(&add)
// 				if _, ok := addList.Addresses[add]; !ok {
// 					continue
// 				}
// 				block := server.RPCBlock{MinedBlock: nil}
// 				ser.MineBlock(add, &block)
// 				if block.MinedBlock == nil {
// 					continue
// 				}
// 				for _, txn := range block.MinedBlock.Txns {
// 					txn.Print()
// 				}
// 			}
// 		case "return":
// 			{
// 				return
// 			}
// 		}
// 	}
// }

// func AddressMenu(ser *server.Server) {
// 	for {
// 		fmt.Println("---Address Menu---")
// 		options := []string{"list", "new", "return"}
// 		option := Menu(options)
// 		switch option {
// 		case "list":
// 			{
// 				addList := server.RPCAddressList{Addresses: map[string]string{}}
// 				ser.GetAddressList(server.Nothing{}, &addList)
// 				for name, pbkey := range addList.Addresses {
// 					fmt.Printf("\t\t %s: %x\n", name, pbkey)
// 				}
// 			}
// 		case "new":
// 			{
// 				addList := server.RPCAddressList{Addresses: map[string]string{}}
// 				ser.GetAddressList(server.Nothing{}, &addList)

// 				fmt.Printf("Enter new address name except ( ")
// 				for name := range addList.Addresses {
// 					fmt.Printf("%s, ", name)
// 				}
// 				fmt.Printf("):")
// 				var add string
// 				fmt.Scan(&add)
// 				status := false
// 				ser.CreateAddress(add, &status)
// 			}
// 		case "return":
// 			{
// 				return
// 			}
// 		}
// 	}
// }

// func TransactionMenu(ser *server.Server) {
// 	for {
// 		fmt.Println("---Transaction Menu---")
// 		options := []string{"new", "return"}
// 		option := Menu(options)
// 		switch option {
// 		case "new":
// 			{
// 				CreateTransaction(ser)
// 			}
// 		case "return":
// 			{
// 				return
// 			}
// 		}
// 	}
// }

// func CreateTransaction(ser *server.Server) {
// 	set := server.RPCTXOSet{Set: map[transaction.OutPoint]transaction.Output{}}
// 	ser.GetUTXO(nil, &set)
// 	txoList := make([]transaction.OutPoint, 0, len(set.Set))
// 	for op := range set.Set {
// 		txoList = append(txoList, op)
// 	}

// 	fmt.Println()
// 	fmt.Println()
// 	for i, out := range txoList {
// 		txo := set.Set[out]
// 		fmt.Printf("%d : %s\n", i, ser.PrintTXO(&out, &txo))
// 	}
// 	fmt.Print("Number of inputs: ")
// 	inSize := 0
// 	fmt.Scan(&inSize)
// 	totalValue := 0
// 	fmt.Print("Enter Inputs\n")
// 	inputs := []transaction.OutPoint{}

// 	for i := 0; i < inSize; i++ {
// 		input := 0
// 		fmt.Scanf("%d", &input)
// 		inputs = append(inputs, txoList[input])
// 		totalValue += int(set.Set[txoList[input]].Value)
// 	}
// 	fmt.Printf("Total Amount : %d\n", totalValue)

// 	addList := server.RPCAddressList{Addresses: map[string]string{}}
// 	ser.GetAddressList(server.Nothing{}, &addList)

// 	fmt.Print("Number of Outputs: ")
// 	outSize := 0
// 	fmt.Scan(&outSize)
// 	totalOutValue := 0
// 	fmt.Print("Enter Outputs\n")
// 	outputs := []transaction.Output{}
// 	for i := 0; i < outSize; i++ {
// 		name := ""
// 		value := 0
// 		fmt.Scanf("%s %d", &name, &value)
// 		address := server.RPCAddress{}
// 		ser.GetAddress(name, &address)
// 		outputs = append(outputs, transaction.Output{Reciever: address.Add, Value: core.Cookies(value)})
// 		totalOutValue += value
// 	}
// 	status := true
// 	ser.SubmitTransaction(&transaction.Transaction{Inputs: inputs, Outputs: outputs}, &status)
// }
