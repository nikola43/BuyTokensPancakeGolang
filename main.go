package main

import (
	"buytokenspancakegolang/models"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	pancakeFactory "buytokenspancakegolang/contracts/IPancakeFactory"
	pancakeRouter "buytokenspancakegolang/contracts/IPancakeRouter02"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hrharder/go-gas"
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
	contractAbi, _ := abi.JSON(strings.NewReader(string(pancakeFactory.PancakeABI)))

	for {
		select {
		case err := <-sub.Err():
			fmt.Println(err)
			//out <- err.Error()

		case vLog := <-logs:
			fmt.Println("vLog.TxHash: " + vLog.TxHash.Hex())
			fmt.Println("vLog.BlockNumber: " + strconv.FormatUint(vLog.BlockNumber, 10))
			res, err := contractAbi.Unpack("PairCreated", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}
			event := new(models.EventsCatched)
			event.TxHash = vLog.TxHash.Hex()
			if res[0].(common.Address) != common.HexToAddress("0xae13d989daC2f0dEbFf460aC112a837C89BAa7cd") {
				event.TokenAddress = res[0].(common.Address)
			} else {
				event.TokenAddress = res[1].(common.Address)
			}
			InsertNewEvent(db, event)
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
		LPAddress:    common.HexToAddress("0"),
		LPPairA:      common.HexToAddress("0"),
		LPPairB:      common.HexToAddress("0"),
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

func Buy(web3GolangHelper *web3helper.Web3GolangHelper, url string) {
	// contract addresses
	pancakeContractAddress := common.HexToAddress("0x10ED43C718714eb63d5aA57B78B54704E256024E") // pancake router address
	wBnbContractAddress := "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"                         // wbnb token adddress
	tokenContractAddress := common.HexToAddress("0xe9C615E0b739e16994a080cA99730Ec104F28CC4")   // eth token adddress

	// create pancakeRouter pancakeRouterInstance
	pancakeRouterInstance, instanceErr := pancakeRouter.NewPancake(pancakeContractAddress, web3GolangHelper.HttpClient())
	if instanceErr != nil {
		fmt.Println(instanceErr)
	}
	fmt.Println("pancakeRouterInstance contract is loaded")

	// calculate gas and gas limit
	gasLimit := uint64(2100000) // in units
	gasPrice, gasPriceErr := gas.SuggestGasPrice(gas.GasPriorityAverage)
	if gasPriceErr != nil {
		fmt.Println(gasPriceErr)
	}

	fmt.Println(

		wBnbContractAddress,
		tokenContractAddress,
		pancakeRouterInstance,
		gasLimit,
		gasPrice,
	)

	/*

		// calculate fee and final value
		gasFee := web3GolangHelper.CalcGasCost(gasLimit, gasPrice)
		ethValue := web3GolangHelper.EtherToWei(big.NewFloat(0.1))
		finalValue := big.NewInt(0).Sub(ethValue, gasFee)

		// set transaction data
		ethBasedClient.ConfigureTransactor(finalValue, gasPrice, gasLimit)
		amountOutMin := big.NewInt(1.0)
		deadline := big.NewInt(time.Now().Unix() + 10000)
		path := ethutils.GeneratePath(wBnbContractAddress, tokenContractAddress.Hex())


		if transactOptsErr {
			fmt.Println(transactOptsErr)
		}

		swapTx, SwapExactETHForTokensErr := pancakeRouterInstance.SwapExactETHForTokensSupportingFeeOnTransferTokens(
			ethBasedClient.Transactor,
			amountOutMin,
			path,
			web3GolangHelper.fromAddress,
			deadline)
		if SwapExactETHForTokensErr != nil {
			fmt.Println("SwapExactETHForTokensErr")
			fmt.Println(SwapExactETHForTokensErr)
		}

		fmt.Println(swapTx)

		txHash := swapTx.Hash().Hex()
		fmt.Println(txHash)
		genericutils.OpenBrowser("https://bscscan.com/tx/" + txHash)

	*/
}
