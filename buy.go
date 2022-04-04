package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/hrharder/go-gas"
	"github.com/joho/godotenv"
	"github.com/nikola43/BuyTokensPancakeGolang/contracts/PancakeRouter"
	"github.com/nikola43/BuyTokensPancakeGolang/errorsutil"
	"github.com/nikola43/BuyTokensPancakeGolang/ethbasedclient"
	"github.com/nikola43/BuyTokensPancakeGolang/ethutils"
	"github.com/nikola43/BuyTokensPancakeGolang/genericutils"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"time"
)

type Wallet struct {
	PublicKey  string `json:"PublicKey"`
	PrivateKey string `json:"PrivateKey"`
}

func main() {

	wallets := make([]Wallet, 0)

	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	wPath := "/home/nkt/wallets"
	files, err := ioutil.ReadDir(wPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fileName := file.Name()
		fmt.Println("fileName", fileName)

		wallet := Wallet{
			PublicKey:  "",
			PrivateKey: "",
		}

		// Open our jsonFile
		jsonFile, _ := os.Open(wPath + "/" + fileName)
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &wallet)
		fmt.Println(wallet)
		wallets = append(wallets, wallet)
	}

	// Open our jsonFile
	jsonFile, err := os.Open("users.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// connect with rpc
	rawurl := "https://bsc-dataseed.binance.org/"

	ethBasedClient := ethbasedclient.New(rawurl, wallets[0].PrivateKey)

	// contract addresses
	pancakeContractAddress := common.HexToAddress("0x10ED43C718714eb63d5aA57B78B54704E256024E") // pancake router address
	wBnbContractAddress := "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"                         // wbnb token adddress

	tokenContractAddress := common.HexToAddress("0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56") // eth token adddress

	// create pancakeRouter pancakeRouterInstance
	pancakeRouterInstance, instanceErr := PancakeRouter.NewPancake(pancakeContractAddress, ethBasedClient.Client)
	errorsutil.HandleError(instanceErr)
	fmt.Println("pancakeRouterInstance contract is loaded")

	// calculate gas and gas limit
	gasLimit := uint64(2100000) // in units
	gasPrice, gasPriceErr := gas.SuggestGasPrice(gas.GasPriorityAverage)
	errorsutil.HandleError(gasPriceErr)

	// calculate fee and final value
	gasFee := ethutils.CalcGasCost(gasLimit, gasPrice)
	ethValue := ethutils.EtherToWei(big.NewFloat(0.01))
	finalValue := big.NewInt(0).Sub(ethValue, gasFee)

	// set transaction data
	ethBasedClient.ConfigureTransactor(finalValue, gasPrice, gasLimit)
	amountOutMin := big.NewInt(1)
	deadline := big.NewInt(time.Now().Unix() + 10000)
	path := ethutils.GeneratePath(wBnbContractAddress, tokenContractAddress.Hex())

	// send transaction
	swapTx, SwapExactETHForTokensErr := pancakeRouterInstance.SwapExactETHForTokensSupportingFeeOnTransferTokens(
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
	genericutils.OpenBrowser("https://bscscan.com/tx/" + txHash)

	tx, err := ethutils.CancelTransaction(ethBasedClient.Client, swapTx, ethBasedClient.PrivateKey)
	errorsutil.HandleError(err)

	txHash = tx.Hash().Hex()
	fmt.Println(txHash)
	genericutils.OpenBrowser("https://bscscan.com/tx/" + txHash)
	os.Exit(0)
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
