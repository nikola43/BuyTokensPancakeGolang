package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/hrharder/go-gas"
	"github.com/joho/godotenv"
	"github.com/nikola43/buy_pancake/contracts/PancakeRouter"
	"github.com/nikola43/buy_pancake/errorsutil"
	"github.com/nikola43/buy_pancake/ethbasedclient"
	"github.com/nikola43/buy_pancake/ethutils"
	"github.com/nikola43/buy_pancake/genericutils"
	"log"
	"math/big"
	"os"
	"time"
)

func main() {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// connect with rpc
	// connect with rpc
	rawurl := "https://bsc-dataseed.binance.org/"
	plainPrivateKey := os.Getenv("PRIVATE_KEY")
	ethBasedClient := ethbasedclient.New(rawurl, plainPrivateKey)

	// contract addresses
	pancakeContractAddress := common.HexToAddress("0x10ed43c718714eb63d5aa57b78b54704e256024e") // pancake router address
	wBnbContractAddress := common.HexToAddress("0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c")    // wbnb token adddress
	tokenContractAddress := common.HexToAddress("0xba2ae424d960c26247dd6c32edc70b295c744c43")   // peg doge token adddress

	// create pancakeRouter pancakeRouterInstance
	pancakeRouterInstance, instanceErr := PancakeRouter.NewPancakeRouter(pancakeContractAddress, ethBasedClient.Client)
	errorsutil.HandleError(instanceErr)
	fmt.Println("pancakeRouterInstance contract is loaded")

	// calculate gas and gas limit
	gasLimit := uint64(210000) // in units
	gasPrice, gasPriceErr := gas.SuggestGasPrice(gas.GasPriorityAverage)
	errorsutil.HandleError(gasPriceErr)

	// calculate fee and final value
	gasFee := ethutils.CalcGasCost(gasLimit, gasPrice)
	ethValue := ethutils.EtherToWei(big.NewFloat(1.0))
	finalValue := big.NewInt(0).Sub(ethValue, gasFee)

	// set transaction data
	ethBasedClient.ConfigureTransactor(finalValue, gasPrice, gasLimit)
	amountOutMin := big.NewInt(1000)
	deadline := big.NewInt(time.Now().Unix() + 10000)
	path := ethutils.GeneratePath(wBnbContractAddress.Hex(), tokenContractAddress.Hex())

	// send transaction
	swapTx, SwapExactETHForTokensErr := pancakeRouterInstance.SwapExactETHForTokens(
		ethBasedClient.Transactor,
		amountOutMin,
		path,
		ethBasedClient.Address,
		deadline)
	if SwapExactETHForTokensErr != nil {
		fmt.Println("SwapExactETHForTokensErr")
		fmt.Println(SwapExactETHForTokensErr)
		os.Exit(0)
	}

	txHash := swapTx.Hash().Hex()
	fmt.Println(txHash)
	genericutils.OpenBrowser("https://bscscan.com/tx/" + txHash)

	/*
		tx, err := ethutils.CancelTransaction(ethBasedClient.Client, swapTx, ethBasedClient.PrivateKey)
		errorsutil.HandleError(err)

		txHash = tx.Hash().Hex()
		fmt.Println(txHash)
		genericutils.OpenBrowser("https://testnet.bscscan.com/tx/" + txHash)
	*/
}

/*
	fmt.Println("ethValue")
	fmt.Println(ethValue)
	fmt.Println("finalValue")
	fmt.Println(finalValue)
	fmt.Println("gasLimit")
	fmt.Println(gasLimit)
	fmt.Println("gasPrice")
	fmt.Println(gasPrice)
	fmt.Println("gasFee")
	fmt.Println(gasFee)
	fmt.Println("nonce")
	fmt.Println(ethBasedClient.Transactor.Nonce)
	fmt.Println("amountOutMin")
	fmt.Println(amountOutMin)
	fmt.Println("path")
	fmt.Println(path)
	fmt.Println("deadline")
	fmt.Println(deadline)
	fmt.Println("transactor")
*/
