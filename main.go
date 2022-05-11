package main

import (
	"buytokenspancakegolang/models"
	"context"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/nikola43/web3golanghelper/web3helper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	pk := "b366406bc0b4883b9b4b3b41117d6c62839174b7d21ec32a5ad0cc76cb3496bd"
	rpcUrl := "https://speedy-nodes-nyc.moralis.io/84a2745d907034e6d388f8d6/bsc/testnet"
	wsUrl := "wss://speedy-nodes-nyc.moralis.io/84a2745d907034e6d388f8d6/bsc/testnet/ws"
	web3GolangHelper := web3helper.NewWeb3GolangHelper(rpcUrl, wsUrl, pk)

	chainID, err := web3GolangHelper.HttpClient().NetworkID(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Chain Id: " + chainID.String())

	db := InitDatabase()
	fmt.Println(db)

	contractAddress := "0xB7926C0430Afb07AA7DEfDE6DA862aE0Bde767bc"
	logs := make(chan types.Log)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(contractAddress)},
	}

	sub, err := web3GolangHelper.WebSocketClient().SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		fmt.Println(sub)
	}

	for {
		select {
		case err := <-sub.Err():
			fmt.Println(err)
			//out <- err.Error()

		case vLog := <-logs:
			fmt.Println("vLog.TxHash: " + vLog.TxHash.Hex())
			fmt.Println("vLog.BlockNumber: " + strconv.FormatUint(vLog.BlockNumber, 10))
			// event := new(models.EventsCatched)
			// event.Tx
			// InsertNewEvent(db)
		}
	}
}

func InitDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func InsertNewEvent(db *gorm.DB, newEvent *models.EventsCatched) bool {
	lpPairs := make([]models.LPPair, 0)
	lpPairs = append(lpPairs, models.LPPair{
		LPAddress:    "",
		LPPairA:      "",
		LPPairB:      "",
		HasLiquidity: false,
	})
	db.Create(&models.EventsCatched{TxHash: newEvent.TxHash, TokenAddress: newEvent.TokenAddress, LPPairs: lpPairs})

	return true
}

func UpdateLiquidity(db *gorm.DB, txHash string) bool {
	var event *models.EventsCatched
	db.First(&event, "TxHash = ?", txHash)

	return true
}

func checkTokens() {
	// Get all records
	//result := db.Find(&users)
	// SELECT * FROM users;
}
