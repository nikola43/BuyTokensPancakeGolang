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
	"github.com/nikola43/buy_pancake/contracts/IERC20"
	"github.com/nikola43/buy_pancake/contracts/PancakeRouter"
	"github.com/nikola43/buy_pancake/errorsutil"
	"github.com/nikola43/buy_pancake/ethutils"
	"github.com/nikola43/buy_pancake/genericutils"
	"log"
	"math/big"
	"os"
	"time"
)

func main() {
	// contract addresses
	// pancakeContractAddress := common.HexToAddress("0x9Ac64Cc6e4415144C455BD8E4837Fea55603e5c3") // pancake router address
	wBnbContractAddress := common.HexToAddress("0xae13d989dac2f0debff460ac112a837c89baa7cd")    // wbnb token adddress
	wBusdContractAddress := common.HexToAddress("0x78867bbeef44f2326bf8ddd1941a4439382ef2a7")   // busd token adddress
	pancakeContractAddress := common.HexToAddress("0x9Ac64Cc6e4415144C455BD8E4837Fea55603e5c3") // pancake router address

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
	IERC20Instance, IERC20InstanceErr := IERC20.NewPancake(wBusdContractAddress, client)
	errorsutil.HandleError(IERC20InstanceErr)
	fmt.Println("IERC20Instance contract is loaded")

	bal, IERC20InstanceBalance := IERC20Instance.BalanceOf(&bind.CallOpts{}, fromAddress)
	errorsutil.HandleError(IERC20InstanceBalance)
	fmt.Println("bal")
	fmt.Println(bal)

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
	transactor.Nonce = big.NewInt(int64(nonce))
	transactor.Context = context.Background()

	fmt.Println("gasLimit")
	fmt.Println(gasLimit)
	fmt.Println("gasPrice")
	fmt.Println(gasPrice)
	fmt.Println("gasFee")
	fmt.Println(gasFee)
	fmt.Println("nonce")
	fmt.Println(nonce)
	fmt.Println("transactor")

	/*
		// create pancakeRouter pancakeRouterInstance
		IPancakePairInstance, IPancakePairErr := IPancakePair.NewPancake(wBusdContractAddress, client)
		errorsutil.HandleError(IPancakePairErr)
		fmt.Println("IPancakePairInstance contract is loaded")
	*/

	swapTx, ApproveErr := IERC20Instance.Approve(
		transactor,
		fromAddress,
		big.NewInt(0).Sub(bal, big.NewInt(0).Div(big.NewInt(0).Mul(bal, big.NewInt(10)), big.NewInt(100))))
	if ApproveErr != nil {
		fmt.Println("ApproveErr")
		fmt.Println(ApproveErr)
	}

	txHash := swapTx.Hash().Hex()
	fmt.Println(txHash)
	genericutils.OpenBrowser("https://testnet.bscscan.com/tx/" + txHash)

	time.Sleep(10 * time.Second)
	fmt.Println("Swapping {tokenValue2} {symbol} for BNB")

	// --------------------------------------------------------------------------------------

	// calculate next nonce
	nonceSell, nonceSellErr := client.PendingNonceAt(context.Background(), fromAddress)
	errorsutil.HandleError(nonceSellErr)

	// calculate gas and gas limit
	gasSellLimit := uint64(210000) // in units

	// get a gas price in base units with one of the exported priorities (fast, fastest, safeLow, average)
	gasSellPrice, gasPriceSellErr := gas.SuggestGasPrice(gas.GasPriorityAverage)
	errorsutil.HandleError(gasPriceSellErr)

	// generate transactor for transactions management
	transactorSell, transactOptsSellErr := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	errorsutil.HandleError(transactOptsSellErr)

	// set transaction data
	transactorSell.GasPrice = gasSellPrice
	transactorSell.GasLimit = gasSellLimit
	gasSellFee := ethutils.CalcGasCost(gasSellLimit, gasSellPrice)
	finalSellValue := big.NewInt(0).Sub(bal, big.NewInt(0).Div(big.NewInt(0).Mul(bal, big.NewInt(10)), big.NewInt(100)))
	transactorSell.Nonce = big.NewInt(int64(nonceSell))
	transactorSell.Context = context.Background()
	deadline := big.NewInt(time.Now().Unix() + 10000)

	path := make([]common.Address, 0)
	path = append(path, wBusdContractAddress)
	path = append(path, wBnbContractAddress)

	fmt.Println("finalSellValueETH")
	fmt.Println(ethutils.WeiToEther(finalSellValue))

	fmt.Println("gasSellFee")
	fmt.Println(gasSellFee)
	fmt.Println("finalSellValue")
	fmt.Println(finalSellValue)
	fmt.Println("gasLimit")
	fmt.Println(gasLimit)
	fmt.Println("gasPrice")
	fmt.Println(gasPrice)
	fmt.Println("gasFee")
	fmt.Println(gasFee)
	fmt.Println("nonce")
	fmt.Println(nonce)
	fmt.Println("path")
	fmt.Println(path)
	fmt.Println("deadline")
	fmt.Println(deadline)
	fmt.Println("transactor")

	// create pancakeRouter pancakeRouterInstance
	pancakeRouterInstance, instanceErr := PancakeRouter.NewPancake(pancakeContractAddress, client)
	errorsutil.HandleError(instanceErr)
	fmt.Println("pancakeRouterInstance contract is loaded")

	swapTx2, SwapExactETHForTokensErr := pancakeRouterInstance.SwapExactTokensForETH(
		transactorSell,
		finalSellValue,
		big.NewInt(0),
		path,
		fromAddress,
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
