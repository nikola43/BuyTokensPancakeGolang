package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/nikola43/web3golanghelper/web3helper"
)

func main() {

	pk := "b366406bc0b4883b9b4b3b41117d6c62839174b7d21ec32a5ad0cc76cb3496bd"
	rpcUrl := "wss://speedy-nodes-nyc.moralis.io/84a2745d907034e6d388f8d6/bsc/testnet/ws"
	web3WsClient := web3helper.NewWsWeb3Client(rpcUrl, pk)

	chainID, err := web3WsClient.NetworkID(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Chain Id: " + chainID.String())

	contractAddress := "0xB7926C0430Afb07AA7DEfDE6DA862aE0Bde767bc"

	//var logs []chan types.Log
	//var subs []ethereum.Subscription

	logs := make(chan types.Log)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(contractAddress)},
	}

	sub, err := web3WsClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		fmt.Println(sub)
	}

	//out := make(chan string)

	for {
		select {
		case err := <-sub.Err():
			fmt.Println(err)
			//out <- err.Error()

		case vLog := <-logs:
			fmt.Println("vLog.TxHash: " + vLog.TxHash.Hex())
			fmt.Println("vLog.BlockNumber: " + strconv.FormatUint(vLog.BlockNumber, 10))
			//out <- vLog.TxHash.Hex()
		}
	}
}