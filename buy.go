package main

import (
	"context"
	"fmt"

	"github.com/nikola43/web3golanghelper/web3helper"
)

func main() {

	pk := "b366406bc0b4883b9b4b3b41117d6c62839174b7d21ec32a5ad0cc76cb3496bd"
	rpcUrl := "wss://speedy-nodes-nyc.moralis.io/84a2745d907034e6d388f8d6/bsc/testnet/ws"
	web3HttpClient := web3helper.NewWsWeb3Client(rpcUrl, pk)

	chainID, err := web3HttpClient.NetworkID(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Chain Id: " + chainID.String())
}