package main

import (
	"buytokenspancakegolang/genericutils"
	"buytokenspancakegolang/models"
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	ierc20 "buytokenspancakegolang/contracts/IERC20"
	pancakeFactory "buytokenspancakegolang/contracts/IPancakeFactory"
	pancakePair "buytokenspancakegolang/contracts/IPancakePair"
	pancakeRouter "buytokenspancakegolang/contracts/IPancakeRouter02"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/fatih/color"
	"github.com/hrharder/go-gas"
	"github.com/mattn/go-colorable"
	"github.com/nikola43/web3golanghelper/web3helper"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Create SprintXxx functions to mix strings with other non-colorized strings:
var yellow = color.New(color.FgYellow).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()

func main() {
	// Declarations
	web3GolangHelper := initWeb3()
	db := InitDatabase()
	factoryAddress := "0xB7926C0430Afb07AA7DEfDE6DA862aE0Bde767bc"
	factoryAbi, _ := abi.JSON(strings.NewReader(string(pancakeFactory.PancakeABI)))

	// LOGIC -----------------------------------------------------------
	checkTokens(db)
	proccessEvents(db, web3GolangHelper, factoryAddress, factoryAbi)
}

func proccessEvents(db *gorm.DB, web3GolangHelper *web3helper.Web3GolangHelper, contractAddress string, contractAbi abi.ABI) {

	wBnbContractAddress := "0xae13d989daC2f0dEbFf460aC112a837C89BAa7cd"
	logs := make(chan types.Log)
	sub := BuildContractEventSubscription(web3GolangHelper, contractAddress, logs)

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
			fmt.Println(res)
			event := new(models.EventsCatched)
			lpPairs := make([]*models.LpPair, 0)
			lpPairs = append(lpPairs, &models.LpPair{
				LPAddress:    common.HexToAddress("0").Hex(),
				LPPairA:      res[0].(common.Address).Hex(),
				LPPairB:      res[1].(common.Address).Hex(),
				HasLiquidity: false,
			})

			event.TxHash = vLog.TxHash.Hex()
			event.LPPairs = lpPairs
			if res[0].(common.Address) != common.HexToAddress(wBnbContractAddress) {
				event.TokenAddress = res[0].(common.Address).Hex()
			} else {
				event.TokenAddress = res[1].(common.Address).Hex()
			}
			InsertNewEvent(db, event)
		}
	}
}

func initWeb3() *web3helper.Web3GolangHelper {
	pk := "b366406bc0b4883b9b4b3b41117d6c62839174b7d21ec32a5ad0cc76cb3496bd"
	rpcUrl := "https://speedy-nodes-nyc.moralis.io/84a2745d907034e6d388f8d6/bsc/testnet"
	wsUrl := "wss://speedy-nodes-nyc.moralis.io/84a2745d907034e6d388f8d6/bsc/testnet/ws"
	web3GolangHelper := web3helper.NewWeb3GolangHelper(rpcUrl, wsUrl, pk)

	chainID, err := web3GolangHelper.HttpClient().NetworkID(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Chain Id: " + chainID.String())
	return web3GolangHelper
}

func BuildContractEventSubscription(web3GolangHelper *web3helper.Web3GolangHelper, contractAddress string, logs chan types.Log) ethereum.Subscription {

	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(contractAddress)},
	}

	sub, err := web3GolangHelper.WebSocketClient().SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		fmt.Println(sub)
	}
	return sub
}

func InitDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func InsertNewEvent(db *gorm.DB, newEvent *models.EventsCatched) bool {
	db.Create(newEvent)

	return true
}

func UpdateLiquidity(db *gorm.DB, eventID uint) bool {
	lpPair := new(models.LpPair)
	db.First(&lpPair, "events_catched_id = ?", eventID).Update("has_liquidity", 1)

	return true
}

func UpdateName(db *gorm.DB, token string, name string) bool {
	event := new(models.EventsCatched)
	db.First(&event, "token_address = ?", token).Update("token_name", 1)

	return true
}

func checkTokens(db *gorm.DB) {
	events := make([]*models.EventsCatched, 0)
	db.Find(&events)
	lo.ForEach(events, func(element *models.EventsCatched, _ int) {
		printTokenStatus(element)
		liquidity := false
		if liquidity {
			UpdateLiquidity(db, element.ID)
		}
	})

}

func Buy(web3GolangHelper *web3helper.Web3GolangHelper, tokenAddress string) {
	// contract addresses
	pancakeContractAddress := common.HexToAddress("0x9Ac64Cc6e4415144C455BD8E4837Fea55603e5c3") // pancake router address
	wBnbContractAddress := "0xae13d989daC2f0dEbFf460aC112a837C89BAa7cd"                         // wbnb token adddress
	tokenContractAddress := common.HexToAddress(tokenAddress)                                   // eth token adddress

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

	// calculate fee and final value
	gasFee := web3helper.CalcGasCost(gasLimit, gasPrice)
	ethValue := web3helper.EtherToWei(big.NewFloat(0.1))
	finalValue := big.NewInt(0).Sub(ethValue, gasFee)

	// set transaction data
	transactor := web3GolangHelper.BuildTransactor(finalValue, gasPrice, gasLimit)
	amountOutMin := big.NewInt(1.0)
	deadline := big.NewInt(time.Now().Unix() + 10000)
	path := web3helper.GeneratePath(wBnbContractAddress, tokenContractAddress.Hex())

	swapTx, SwapExactETHForTokensErr := pancakeRouterInstance.SwapExactETHForTokensSupportingFeeOnTransferTokens(
		transactor,
		amountOutMin,
		path,
		*web3GolangHelper.FromAddress,
		deadline)
	if SwapExactETHForTokensErr != nil {
		fmt.Println("SwapExactETHForTokensErr")
		fmt.Println(SwapExactETHForTokensErr)
	}

	fmt.Println(swapTx)

	txHash := swapTx.Hash().Hex()
	fmt.Println(txHash)
	genericutils.OpenBrowser("https://testnet.bscscan.com/tx/" + txHash)

}

/*
   function swapExactETHForTokensSupportingFeeOnTransferTokens(
       uint amountOutMin,
       address[] calldata path,
       address to,
       uint deadline
   ) external payable;
*/

func BuyV2(web3GolangHelper *web3helper.Web3GolangHelper, tokenAddress string, value *big.Int) {
	toAddress := common.HexToAddress("0x9Ac64Cc6e4415144C455BD8E4837Fea55603e5c3")
	wBnbContractAddress := "0xae13d989daC2f0dEbFf460aC112a837C89BAa7cd"

	transferFnSignature := []byte("swapExactETHForTokensSupportingFeeOnTransferTokens(uint,address[],address,uint)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	path := web3helper.GeneratePath(wBnbContractAddress, tokenAddress)
	pathString := []string{path[0].Hex(), path[1].Hex()}

	deadline := big.NewInt(time.Now().Unix() + 10000)
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(pathString)
	bs := buf.Bytes()
	fmt.Printf("%q", bs)

	paddedAmountOutMin := common.LeftPadBytes(value.Bytes(), 32)
	paddedPathA := common.LeftPadBytes(path[0].Bytes(), 32)
	paddedPathB := common.LeftPadBytes(path[1].Bytes(), 32)
	paddedPath := common.LeftPadBytes(bs, 32)
	paddedTo := common.LeftPadBytes(toAddress.Bytes(), 32)
	paddedDeadline := common.LeftPadBytes(deadline.Bytes(), 32)

	fmt.Println("paddedAmountOutMin", paddedAmountOutMin)
	fmt.Println("paddedPathA", paddedPathA)
	fmt.Println("paddedPathB", paddedPathB)
	fmt.Println("paddedPath", paddedPath)
	fmt.Println("paddedTo", paddedTo)
	fmt.Println("paddedDeadline", paddedDeadline)
	fmt.Println("paddedAmountOutMin", paddedAmountOutMin)
	fmt.Println("paddedAmountOutMin", paddedAmountOutMin)

	txData := web3helper.BuildTxData(methodID, paddedAmountOutMin, paddedPath, paddedTo, paddedDeadline)

	fmt.Println("txData", txData)

	estimateGas := web3GolangHelper.EstimateGas(toAddress.Hex(), txData)

	fmt.Println("estimateGas", estimateGas)

	txId, txNonce, err := web3GolangHelper.SignAndSendTransaction(toAddress.Hex(), web3helper.ToWei(value, 18), txData, web3GolangHelper.PendingNonce(), nil, estimateGas)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(txId)
	fmt.Println(txNonce)
}

/*
	toAddress := common.HexToAddress(toAddressString)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(value.Bytes(), 32)

	txData := BuildTxData(methodID, paddedAddress, paddedAmount)

	estimateGas := w.EstimateGas(tokenAddressString, txData)
	txId, txNonce, err := w.SignAndSendTransaction(toAddressString, ToWei(value, 18), txData, w.PendingNonce(), nil, estimateGas)
	if err != nil {
		return "", big.NewInt(0), err
	}

	return txId, txNonce, nil
*/

func updateTokenStatus(db *gorm.DB, web3GolangHelper *web3helper.Web3GolangHelper, token *models.EventsCatched) {

	// create pancakeRouter pancakeRouterInstance
	tokenContractInstance, instanceErr := ierc20.NewPancake(common.HexToAddress(token.TokenAddress), web3GolangHelper.HttpClient())
	if instanceErr != nil {
		fmt.Println(instanceErr)
	}
	fmt.Println("pancakeRouterInstance contract is loaded")

	tokenName, getNameErr := tokenContractInstance.Name(nil)
	if getNameErr != nil {
		UpdateName(db, token.TokenAddress, tokenName)
		fmt.Println(getNameErr)
	}

	fmt.Println(tokenName)

}

func getTokenPairs(web3GolangHelper *web3helper.Web3GolangHelper, token *models.EventsCatched) string {
	//lpPairs := make([]*models.LpPair, 0)

	lpPairAddress := getPair(web3GolangHelper, token.TokenAddress)

	//append(lpPairs, )

	fmt.Println("lpPairAddress", lpPairAddress)
	return lpPairAddress
}

func getReserves(web3GolangHelper *web3helper.Web3GolangHelper, tokenAddress string) struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
} {

	pairInstance, instanceErr := pancakePair.NewPancake(common.HexToAddress("0xB7926C0430Afb07AA7DEfDE6DA862aE0Bde767bc"), web3GolangHelper.HttpClient())
	if instanceErr != nil {
		fmt.Println(instanceErr)
	}

	reserves, getReservesErr := pairInstance.GetReserves(nil)
	if getReservesErr != nil {
		fmt.Println(getReservesErr)
	}

	return reserves
}

func getPair(web3GolangHelper *web3helper.Web3GolangHelper, tokenAddress string) string {

	factoryInstance, instanceErr := pancakeFactory.NewPancake(common.HexToAddress("0xB7926C0430Afb07AA7DEfDE6DA862aE0Bde767bc"), web3GolangHelper.HttpClient())
	if instanceErr != nil {
		fmt.Println(instanceErr)
	}

	wBnbContractAddress := "0xae13d989daC2f0dEbFf460aC112a837C89BAa7cd"

	lpPairAddress, getPairErr := factoryInstance.GetPair(nil, common.HexToAddress(wBnbContractAddress), common.HexToAddress(tokenAddress))
	if getPairErr != nil {
		fmt.Println(getPairErr)
	}

	return lpPairAddress.Hex()

}

func printTokenStatus(token *models.EventsCatched) {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.Info("TOKEN INFO")

	fmt.Printf("%s: %s\n", cyan("Token Address"), yellow(token.TokenAddress))
	fmt.Printf("%s:\n", cyan("LP Pairs"))
	lo.ForEach(token.LPPairs, func(element *models.LpPair, _ int) {
		fmt.Printf("\t%s: %s\n", cyan("LP Address"), yellow(element.LPAddress))
		fmt.Printf("\t%s: %s\n", cyan("LP TokenA Address"), yellow(element.LPPairA))
		fmt.Printf("\t%s: %s\n", cyan("LP TokenB Address"), yellow(element.LPPairB))
		fmt.Printf("\t%s: %s\n\n", cyan("LP Has Liquidity"), getPairLiquidityIcon(element))
		fmt.Printf("\t%s: %s\n\n", cyan("Trading Enabled"), getPairTradingIcon(element))
	})
}

func getPairTradingIcon(pair *models.LpPair) string {
	icon := "🔴"
	if pair.TradingEnabled {
		icon = "🟢"
	}
	return icon
}

func getPairLiquidityIcon(pair *models.LpPair) string {
	icon := "🔴"
	if pair.HasLiquidity {
		icon = "🟢"
	}
	return icon
}
