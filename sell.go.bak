package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
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
	rawurl := "https://data-seed-prebsc-1-s1.binance.org:8545/"
	plainPrivateKey := os.Getenv("PRIVATE_KEY")
	ethBasedClient := ethbasedclient.New(rawurl, plainPrivateKey)

	// contract addresses
	pancakeContractAddress := common.HexToAddress("0x9Ac64Cc6e4415144C455BD8E4837Fea55603e5c3") // pancake router address
	wBnbContractAddress := common.HexToAddress("0xae13d989dac2f0debff460ac112a837c89baa7cd")    // wbnb token adddress
	tokenContractAddress := common.HexToAddress("0x8babbb98678facc7342735486c851abd7a0d17ca")   // eth token adddress

	// connect with bsc
	client, err := ethclient.Dial("https://data-seed-prebsc-1-s1.binance.org:8545/")
	errorsutil.HandleError(err)
	defer client.Close()

	// create pancakeRouter pancakeRouterInstance
	IERC20Instance, IERC20InstanceErr := IERC20.NewPancake(tokenContractAddress, client)
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
	genericutils.OpenBrowser("https://testnet.bscscan.com/tx/" + txHash)

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
	pancakeRouterInstance, instanceErr := PancakeRouter.NewPancake(pancakeContractAddress, client)
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
	genericutils.OpenBrowser("https://testnet.bscscan.com/tx/" + txHash2)

	os.Exit(0)
}
