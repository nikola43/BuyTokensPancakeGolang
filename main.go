package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hrharder/go-gas"
	"github.com/nikola43/buy_pancake/contracts/PancakeRouter"
	"github.com/nikola43/buy_pancake/errorsutil"
	"github.com/nikola43/buy_pancake/ethutils"
	"time"

	"log"
	"math/big"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	// contract addresses
	pancakeContractAddress := common.HexToAddress("0x9Ac64Cc6e4415144C455BD8E4837Fea55603e5c3") // pancake router address
	wBnbContractAddress := common.HexToAddress("0xae13d989dac2f0debff460ac112a837c89baa7cd")    // wbnb token adddress
	wBusdContractAddress := common.HexToAddress("0x78867bbeef44f2326bf8ddd1941a4439382ef2a7")   // busd token adddress

	// connect with bsc
	client, err := ethclient.Dial("https://data-seed-prebsc-1-s1.binance.org:8545/")
	errorsutil.HandleError(err)
	defer client.Close()

	// create privateKey from string key
	privateKey, privateKeyErr := crypto.HexToECDSA("")
	errorsutil.HandleError(privateKeyErr)

	// generate public key and address from private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	// generate address from public key
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// get chain id
	chainID, chainIDErr := client.ChainID(context.Background())
	errorsutil.HandleError(chainIDErr)

	// get current balance
	balance, balanceErr := client.BalanceAt(context.Background(), fromAddress, nil)
	errorsutil.HandleError(balanceErr)
	fmt.Println(balance) // 25893180161173005034

	// generate transactor for transactions management
	transactor, transactOptsErr := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	errorsutil.HandleError(transactOptsErr)

	// create pancakeRouter pancakeRouterInstance
	pancakeRouterInstance, instanceErr := PancakeRouter.NewPancake(pancakeContractAddress, client)
	errorsutil.HandleError(instanceErr)
	fmt.Println("contract is loaded")

	// calculate next nonce
	nonce, nonceErr := client.PendingNonceAt(context.Background(), fromAddress)
	errorsutil.HandleError(nonceErr)

	// calculate gas and gas limit
	gasLimit := uint64(210000) // in units

	// get a gas price in base units with one of the exported priorities (fast, fastest, safeLow, average)
	gasPrice, gasPriceErr := gas.SuggestGasPrice(gas.GasPriorityAverage)
	errorsutil.HandleError(gasPriceErr)

	// set transaction data
	transactor.GasPrice = gasPrice
	transactor.GasLimit = gasLimit
	gasFee := ethutils.CalcGasCost(gasLimit, gasPrice)
	ethValue := ethutils.EtherToWei(big.NewFloat(1.0))
	finalValue := big.NewInt(0).Sub(ethValue, gasFee)
	transactor.Nonce = big.NewInt(int64(nonce))
	transactor.Value = finalValue
	transactor.Context = context.Background()
	amountOutMin := big.NewInt(1000)
	deadline := big.NewInt(time.Now().Unix() + 10000)

	path := make([]common.Address, 0)
	path = append(path, wBnbContractAddress)
	path = append(path, wBusdContractAddress)

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
	fmt.Println(nonce)
	fmt.Println("amountOutMin")
	fmt.Println(amountOutMin)
	fmt.Println("path")
	fmt.Println(path)
	fmt.Println("deadline")
	fmt.Println(deadline)
	fmt.Println("transactor")


	swapTx, SwapExactETHForTokensErr := pancakeRouterInstance.SwapExactETHForTokens(
		transactor,
		amountOutMin,
		path,
		fromAddress,
		deadline)
	if SwapExactETHForTokensErr != nil {
		fmt.Println("SwapExactETHForTokensErr")
		fmt.Println(SwapExactETHForTokensErr)
	}

	txHash := swapTx.Hash().Hex()
	fmt.Println(txHash)
	openBrowser("https://testnet.bscscan.com/tx/" + txHash)
	os.Exit(0)
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
