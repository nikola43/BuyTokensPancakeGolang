package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/hrharder/go-gas"
	"github.com/joho/godotenv"
	"github.com/nikola43/buy_pancake/contracts/PancakeRouter"
	"github.com/nikola43/buy_pancake/ethbasedclient"
	"github.com/nikola43/buy_pancake/utils/errorsutil"
	"github.com/nikola43/buy_pancake/utils/ethutil"
	"github.com/nikola43/buy_pancake/utils/genericutil"
	"math/big"
	"os"
	"time"
)

func main() {
	// load .env file
	err := godotenv.Load(".env")
	errorsutil.HandleError(err)

	// connect with rpc
	rawurl := "https://data-seed-prebsc-1-s1.binance.org:8545/"
	plainPrivateKey := os.Getenv("PRIVATE_KEY")
	ethBasedClient := ethbasedclient.New(rawurl, plainPrivateKey)

	// contract addresses
	pancakeContractAddress := common.HexToAddress("0x9Ac64Cc6e4415144C455BD8E4837Fea55603e5c3") // pancake router address
	wBnbContractAddress := common.HexToAddress("0xae13d989dac2f0debff460ac112a837c89baa7cd")    // wbnb token adddress
	tokenContractAddress := common.HexToAddress("0x8babbb98678facc7342735486c851abd7a0d17ca")   // eth token adddress

	// create pancakeRouter pancakeRouterInstance
	pancakeRouterInstance, instanceErr := PancakeRouter.NewPancake(pancakeContractAddress, ethBasedClient.Client)
	errorsutil.HandleError(instanceErr)
	fmt.Println("pancakeRouterInstance contract is loaded")

	// calculate gas and gas limit
	gasLimit := uint64(210000) // in units
	gasPrice, gasPriceErr := gas.SuggestGasPrice(gas.GasPriorityAverage)
	errorsutil.HandleError(gasPriceErr)

	// calculate fee and final value
	gasFee := ethutil.CalcGasCost(gasLimit, gasPrice)
	ethValue := ethutil.EtherToWei(big.NewFloat(1.0))
	finalValue := big.NewInt(0).Sub(ethValue, gasFee)

	// set transaction data
	ethBasedClient.ConfigureTransactor(finalValue, gasPrice, gasLimit)
	amountOutMin := big.NewInt(1000)
	deadline := big.NewInt(time.Now().Unix() + 10000)
	path := ethutil.GeneratePath(wBnbContractAddress.Hex(), tokenContractAddress.Hex())

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
	}

	txHash := swapTx.Hash().Hex()
	fmt.Println(txHash)
	genericutil.OpenBrowser("https://testnet.bscscan.com/tx/" + txHash)

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
