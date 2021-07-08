package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nikola43/buy_pancake/contracts/PancakeRouter"
	"log"
	"math/big"
	"os"
)

func main() {

	parseEthToWei(1)




	// contract addresses
	pancakeContractAddress := common.HexToAddress("0x9Ac64Cc6e4415144C455BD8E4837Fea55603e5c3") // pancake router address
	wBnbContractAddress := common.HexToAddress("0xae13d989dac2f0debff460ac112a837c89baa7cd")    // wbnb token adddress
	wBusdContractAddress := common.HexToAddress("0x78867bbeef44f2326bf8ddd1941a4439382ef2a7")   // busd token adddress

	// connect with bsc
	client, err := ethclient.Dial("https://data-seed-prebsc-1-s1.binance.org:8545/")
	handleError(err)
	defer client.Close()

	// create privateKey from string key
	privateKey, privateKeyErr := crypto.HexToECDSA("")
	handleError(privateKeyErr)

	// generate public key and address from private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// get chain id
	chainID, chainIDErr := client.ChainID(context.Background())
	handleError(chainIDErr)

	// get current balance
	balance, balanceErr := client.BalanceAt(context.Background(), fromAddress, nil)
	handleError(balanceErr)
	fmt.Println(balance) // 25893180161173005034

	// generate transactor for transactions management
	transactor, transactOptsErr := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	handleError(transactOptsErr)

	// create pancakeRouter pancakeRouterInstance
	pancakeRouterInstance, instanceErr := PancakeRouter.NewPancake(pancakeContractAddress, client)
	handleError(instanceErr)
	fmt.Println("contract is loaded")



	nonce, nonceErr := client.PendingNonceAt(context.Background(), fromAddress)
	if nonceErr != nil {
		log.Fatal(nonceErr)
	}

	// calculate gas and gas limit
	gasLimit := uint64(210000) // in units
	gasPrice, gasPrice2Err := client.SuggestGasPrice(context.Background())
	if gasPrice2Err != nil {
		log.Fatal(gasPrice2Err)
	}

	transactor.Nonce = big.NewInt(int64(nonce))
	transactor.Value = big.NewInt(10000000000000)
	transactor.GasPrice = gasPrice
	transactor.GasLimit = gasLimit
	transactor.Context = context.Background()
	amountOutMin := big.NewInt(1000)
	deadline := big.NewInt(1625784522)

	path := make([]common.Address, 0)
	path = append(path, wBnbContractAddress)
	path = append(path, wBusdContractAddress)

	fmt.Println(transactor)
	fmt.Println(amountOutMin)
	fmt.Println(path)
	fmt.Println(deadline)

	swapTx, SwapExactETHForTokensErr := pancakeRouterInstance.SwapExactETHForTokens(
		transactor,
		amountOutMin,
		path,
		fromAddress,
		deadline)
	if SwapExactETHForTokensErr != nil {
		fmt.Println("errr")
		fmt.Println(SwapExactETHForTokensErr)
	}

	fmt.Println(swapTx.Hash().Hex())
	os.Exit(0)

}

func parseEthToWei(value float64) int64 {
	fmt.Println(value / 1000000000000000000)

	fmt.Println()

	return int64(value / 1000000000000000000)
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
