package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/fatih/color"
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

	//generateContractEventSubscription
	logs := make(chan types.Log)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(contractAddress)},
	}

	sub, err := w.wsClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return nil, nil, err
	}

	contractLog, contractSub, err := w.GenerateContractEventSubscription(contractsAddresses[i])
	if err != nil {
		return err
	}

	logs = append(logs, contractLog)
	subs = append(subs, contractSub)

	contractAbi, err := abi.JSON(strings.NewReader(string(NodeManagerV83.NodeManagerV83ABI)))
	if err != nil {
		log.Fatal(err)
	}

	go func(i int) {
		fmt.Println(color.MagentaString("    Init Subscription: "), color.YellowString(contractsAddresses[i]))

		for {
			select {
			case err := <-subs[i].Err():
				out <- err.Error()

			case vLog := <-logs[i]:
				//fmt.Println(vLog) // pointer to event log
				// fmt.Println("Data logs")
				//fmt.Println(string(vLog.Data))
				// //fmt.Println("vLog.Address: " + vLog.Address.Hex())
				// fmt.Println("vLog.TxHash: " + vLog.TxHash.Hex())
				// fmt.Println("vLog.BlockNumber: " + strconv.FormatUint(vLog.BlockNumber, 10))
				// fmt.Println("")

				res, err := contractAbi.Unpack("GiftCardPayed", vLog.Data)
				if err != nil {
					log.Fatal(err)
				}
				//fmt.Println("vLog.TxHash: " + vLog.TxHash.Hex())
				//fmt.Println(res) // foo
				services.GetGiftCardIntentPayment(res[2].(string))
				for i := range res {
					fmt.Println(res[i])
				}
				//var topics [4]string
				//fmt.Println(vLog)
				// for i := range vLog.Topics {
				// 	//topics[i] = vLog.Topics[i].Hex()
				// 	fmt.Println(vLog.Topics[i].Hex())
				// }
				fmt.Println("end for loop")
				out <- vLog.TxHash.Hex()
			}
		}
	}(i)
}
