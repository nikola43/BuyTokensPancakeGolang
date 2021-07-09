package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/hrharder/go-gas"
	"github.com/joho/godotenv"
	"github.com/nikola43/buy_pancake/contracts/IERC20"
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
	rawurl := "https://bsc-dataseed.binance.org/"
	plainPrivateKey := os.Getenv("PRIVATE_KEY")
	ethBasedClient := ethbasedclient.New(rawurl, plainPrivateKey)

	// contract addresses
	pancakeContractAddress := common.HexToAddress("0x10ed43c718714eb63d5aa57b78b54704e256024e") // pancake router address
	wBnbContractAddress := common.HexToAddress("0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c")    // wbnb token adddress
	tokenContractAddress := common.HexToAddress("0xba2ae424d960c26247dd6c32edc70b295c744c43")   // peg doge token adddress

	// create pancakeRouter pancakeRouterInstance
	IERC20Instance, IERC20InstanceErr := IERC20.NewIERC20(tokenContractAddress, ethBasedClient.Client)
	errorsutil.HandleError(IERC20InstanceErr)
	fmt.Println("IERC20Instance contract is loaded")

	bal, IERC20InstanceBalance := IERC20Instance.BalanceOf(&bind.CallOpts{}, ethBasedClient.Address)
	errorsutil.HandleError(IERC20InstanceBalance)
	fmt.Println("bal")
	fmt.Println(bal)

	// calculate gas and gas limit
	gasLimit := uint64(210000) // in units
	gasPrice, gasPriceErr := gas.SuggestGasPrice(gas.GasPriorityAverage)
	gasFee := ethutils.CalcGasCost(gasLimit, gasPrice)
	fmt.Println(gasFee)
	errorsutil.HandleError(gasPriceErr)

	// set transaction data
	ethBasedClient.ConfigureTransactor(big.NewInt(-1), gasPrice, gasLimit)
	deadline := big.NewInt(time.Now().Unix() + 10000)

	swapTx, ApproveErr := IERC20Instance.Approve(
		ethBasedClient.Transactor,
		ethBasedClient.Address,
		big.NewInt(0).Sub(bal, big.NewInt(0).Div(big.NewInt(0).Mul(bal, big.NewInt(10)), big.NewInt(100))))
	if ApproveErr != nil {
		fmt.Println("ApproveErr")
		fmt.Println(ApproveErr)
	}

	txHash := swapTx.Hash().Hex()
	fmt.Println(txHash)
	genericutils.OpenBrowser("https://bscscan.com/tx/" + txHash)

	time.Sleep(1 * time.Second)
	fmt.Println("Swapping {tokenValue2} {symbol} for BNB")

	// --------------------------------------------------------------------------------------

	// calculate gas and gas limit
	gasSellLimit := uint64(210000) // in units
	gasSellPrice, gasPriceSellErr := gas.SuggestGasPrice(gas.GasPriorityAverage)
	errorsutil.HandleError(gasPriceSellErr)

	// set transaction data
	ethBasedClient.ConfigureTransactor(big.NewInt(-1), gasPrice, gasLimit)
	deadline = big.NewInt(time.Now().Unix() + 10000)
	path := ethutils.GeneratePath(tokenContractAddress.Hex(), wBnbContractAddress.Hex())

	gasSellFee := ethutils.CalcGasCost(gasSellLimit, gasSellPrice)
	finalSellValue := big.NewInt(0).Sub(bal, big.NewInt(0).Div(big.NewInt(0).Mul(bal, big.NewInt(10)), big.NewInt(100)))

	fmt.Println(gasSellFee)

	// create pancakeRouter pancakeRouterInstance
	pancakeRouterInstance, instanceErr := PancakeRouter.NewPancakeRouter(pancakeContractAddress, ethBasedClient.Client)
	errorsutil.HandleError(instanceErr)
	fmt.Println("pancakeRouterInstance contract is loaded")

	swapTx2, SwapExactETHForTokensErr := pancakeRouterInstance.SwapExactTokensForETH(
		ethBasedClient.Transactor,
		finalSellValue,
		big.NewInt(0),
		path,
		ethBasedClient.Address,
		deadline)
	if SwapExactETHForTokensErr != nil {
		fmt.Println("SwapExactETHForTokensErr")
		fmt.Println(SwapExactETHForTokensErr)
	}

	txHash2 := swapTx2.Hash().Hex()
	fmt.Println(txHash)
	genericutils.OpenBrowser("https://bscscan.com/tx/" + txHash2)

	os.Exit(0)
}
