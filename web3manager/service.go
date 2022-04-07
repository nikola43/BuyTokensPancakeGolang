package web3manager

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/nikola43/BuyTokensPancakeGolang/errorsutil"
	"github.com/nikola43/BuyTokensPancakeGolang/ethutils"
	web3util "github.com/nikola43/BuyTokensPancakeGolang/web3manager/util"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"
	"net/url"
	"strconv"
	"time"
)

type LogLevel int

const (
	NoneLogLevel   LogLevel = 0
	LowLogLevel    LogLevel = 1
	MediumLogLevel LogLevel = 2
	HighLogLevel   LogLevel = 3
)

var chainId = big.NewInt(43113)
var defaultGasLimit = uint64(7000000)
var logLevel = HighLogLevel

type Web3Manager struct {
	plainPrivateKey string
	httpClient      *ethclient.Client
	wsClient        *ethclient.Client
	fromAddress     *common.Address
}

func (w *Web3Manager) AddHttpClient(httpClient *ethclient.Client) error {

	if w.httpClient != nil {
		return errors.New("web3 Http provider already instanced")
	}

	w.httpClient = httpClient
	return nil
}

func (w *Web3Manager) AddWsClient(wsClient *ethclient.Client) error {

	if w.wsClient != nil {
		return errors.New("web3 websocket provider already instanced")
	}

	w.wsClient = wsClient
	return nil
}

func (w *Web3Manager) GenerateContractEventSubscription(contractAddress string) (chan types.Log, ethereum.Subscription, error) {

	logs := make(chan types.Log)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(contractAddress)},
	}

	sub, err := w.wsClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return nil, nil, err
	}

	return logs, sub, nil
}

func (w *Web3Manager) SuggestGasPrice() *big.Int {

	gasPrice, err := w.selectClient().SuggestGasPrice(context.Background())

	if err != nil {
		fmt.Println(err)
		return big.NewInt(0)
	}

	return gasPrice
}

func (w *Web3Manager) ListenBridgesEventsV2(contractsAddresses []string, out chan<- string) error {

	var logs []chan types.Log
	var subs []ethereum.Subscription

	fmt.Println("")
	fmt.Println(color.YellowString("  --------------------- Contracts Subscriptions ---------------------"))
	for i := 0; i < len(contractsAddresses); i++ {

		contractLog, contractSub, err := w.GenerateContractEventSubscription(contractsAddresses[i])
		if err != nil {
			return err
		}

		logs = append(logs, contractLog)
		subs = append(subs, contractSub)

		go func(i int) {
			fmt.Println(color.MagentaString("    Init Subscription: "), color.YellowString(contractsAddresses[i]))

			for {
				select {
				case err := <-subs[i].Err():
					out <- err.Error()

				case vLog := <-logs[i]:
					//fmt.Println(vLog) // pointer to event log
					fmt.Println("Data logs")
					fmt.Println(string(vLog.Data))
					//fmt.Println("vLog.Address: " + vLog.Address.Hex())
					fmt.Println("vLog.TxHash: " + vLog.TxHash.Hex())
					fmt.Println("vLog.BlockNumber: " + strconv.FormatUint(vLog.BlockNumber, 10))
					fmt.Println("")
					out <- vLog.TxHash.Hex()
				}
			}
		}(i)
	}
	return nil
}

func NewHttpWeb3Client(rpcUrl string, plainPrivateKey interface{}) *Web3Manager {

	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	_, getBlockErr := client.BlockNumber(context.Background())
	if getBlockErr != nil {
		log.Fatal(getBlockErr)
	}

	web3Manager := &Web3Manager{
		httpClient: client,
	}

	if plainPrivateKey != nil {
		web3Manager.plainPrivateKey = plainPrivateKey.(string)
		web3Manager.fromAddress = GeneratePublicAddressFromPrivateKey(web3Manager.plainPrivateKey)
	}

	return web3Manager
}

func (w *Web3Manager) CurrentBlockNumber() uint64 {

	blockNumber, getBlockErr := w.selectClient().BlockNumber(context.Background())
	if getBlockErr != nil {
		fmt.Println(getBlockErr)
		return 0
	}

	return blockNumber
}

func (w *Web3Manager) HttpClient() *ethclient.Client {
	return w.httpClient
}

func (w *Web3Manager) WebSocketClient() *ethclient.Client {
	return w.wsClient
}

func (w *Web3Manager) SetPrivateKey(plainPrivateKey string) *Web3Manager {
	w.plainPrivateKey = plainPrivateKey
	return w
}

func NewWsWeb3Client(rpcUrl string, plainPrivateKey interface{}) *Web3Manager {

	_, err := url.ParseRequestURI(rpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	wsClient, err2 := ethclient.Dial(rpcUrl)
	if err2 != nil {
		log.Fatal(err)
	}

	_, getBlockErr := wsClient.BlockNumber(context.Background())
	if getBlockErr != nil {
		log.Fatal(err)
	}

	web3Manager := &Web3Manager{
		wsClient: wsClient,
	}

	if plainPrivateKey != nil {
		web3Manager.plainPrivateKey = plainPrivateKey.(string)
		web3Manager.fromAddress = web3util.GeneratePublicAddressFromPrivateKey(web3Manager.plainPrivateKey)
	}

	return web3Manager
}

func (w *Web3Manager) Unsubscribe() {
	time.Sleep(10 * time.Second)
	fmt.Println("---unsubscribe-----")
	//w.ethSubscription.Unsubscribe()
}

func (w *Web3Manager) GetEthBalance(address string) *big.Int {
	account := common.HexToAddress(address)
	balance, err := w.httpClient.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil
	}
	return balance
}

func (w *Web3Manager) IsAddressContract(address string) bool {

	if !web3util.ValidateAddress(address) {
		return false
	}

	bytecode, err := w.httpClient.CodeAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		return false
	}
	return len(bytecode) > 0
}

func (w *Web3Manager) ChainId() *big.Int {
	chainID, err := w.httpClient.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return chainID
}

func (w *Web3Manager) PendingNonce() uint64 {

	nonce, err := w.selectClient().PendingNonceAt(context.Background(), *w.fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	return nonce
}
func (w *Web3Manager) SignTx(tx *types.Transaction) (*types.Transaction, error) {

	privateKey, privateKeyErr := crypto.HexToECDSA(w.plainPrivateKey)
	if privateKeyErr != nil {
		return nil, privateKeyErr
	}

	signedTx, signTxErr := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if signTxErr != nil {
		return nil, signTxErr
	}

	return signedTx, nil
}

func (w *Web3Manager) NewContract(contractAddress string) {

	/*
		address := common.HexToAddress(contractAddress)
		instance, err := store.NewStore(address, w.httpClient)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("contract is loaded")
		return instance
	*/
}

func (w *Web3Manager) SubscribeContractBridgeBSCEvent(contractAddressString string) error {

	if w.wsClient == nil {
		return errors.New("Nil Web3 Websocket Client")
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(contractAddressString)},
	}

	logs := make(chan types.Log)
	sub, err := w.wsClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Init Sub")
	for {
		select {
		case err := <-sub.Err():
			fmt.Println("Error")
			fmt.Println(err)
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println("Data")
			fmt.Println(string(vLog.Data))
			//fmt.Println("vLog.Address: " + vLog.Address.Hex())
			fmt.Println("vLog.TxHash: " + vLog.TxHash.Hex())
			fmt.Println("vLog.BlockNumber: " + strconv.FormatUint(vLog.BlockNumber, 10))

			/*

					event := struct {
						Key   [32]byte
						Value [32]byte
					}{}


				contractAbi, err := abi.JSON(strings.NewReader(bridgeAvax.BridgeAvaxMetaData.ABI))
				if err != nil {
					log.Fatal(err)
				}

				//r, err := contractAbi.Unpack(&event, "ItemSet", vLog.Data)
				r, err := contractAbi.Unpack("Transfer", vLog.Data)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(r)

			*/

			//fmt.Println(string(event.Key[:]))   // foo
			//fmt.Println(string(event.Value[:])) // bar

			fmt.Println("")
			//fmt.Println(vLog) // pointer to event log
		}
	}
}

func (w *Web3Manager) EstimateGas(to string, txData []byte) uint64 {
	toAddress := common.HexToAddress(to)
	estimateGas, estimateGasErr := w.selectClient().EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: txData,
	})
	if estimateGasErr != nil {
		return 0
	}
	return estimateGas
}

func (w *Web3Manager) SendTokens(tokenAddressString, toAddressString string, value *big.Int) (string, uint64, error) {

	toAddress := common.HexToAddress(toAddressString)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(value.Bytes(), 32)

	txData := web3util.BuildTxData(methodID, paddedAddress, paddedAmount)

	estimateGas := w.EstimateGas(tokenAddressString, txData)
	txId, txNonce, err := w.SignAndSendTransaction(toAddressString, web3util.ToWei(value, 18), txData, w.PendingNonce(), nil, estimateGas)
	if err != nil {
		return "", 0, err
	}

	return txId, txNonce, nil
}

func (w *Web3Manager) selectClient() *ethclient.Client {
	var selectedClient *ethclient.Client
	if w.wsClient != nil {
		selectedClient = w.wsClient
	} else {
		if w.httpClient != nil {
			selectedClient = w.httpClient
		} else {
			log.Fatal("SuggestGasPrice: Not conected")
		}
	}
	return selectedClient
}

func (w *Web3Manager) SendEth(toAddressString string, value string) (string, uint64, error) {

	txId, nonce, err := w.SignAndSendTransaction(toAddressString, web3util.ToWei(value, 18), make([]byte, 0), w.PendingNonce(), nil, nil)
	if err != nil {
		return "", 0, err
	}

	return txId, nonce, nil
}

func (w *Web3Manager) SignAndSendTransaction(toAddressString string, value *big.Int, data []byte, nonce uint64, customGasPrice interface{}, customGasLimit interface{}) (string, uint64, error) {

	usedGasPrice, _ := w.selectClient().SuggestGasPrice(context.Background())
	if logLevel == MediumLogLevel {
		fmt.Println(color.CyanString("usedGasPrice -> suggestGasPrice: "), color.YellowString(strconv.Itoa(int(usedGasPrice.Int64())))+"\n")
	}

	if customGasPrice != nil {
		usedGasPrice = customGasPrice.(*big.Int)

		if logLevel == MediumLogLevel {
			fmt.Println(color.CyanString("usedGasPrice -> customGasPrice: "), color.YellowString(strconv.Itoa(int(usedGasPrice.Int64())))+"\n")
		}
	}

	usedGasLimit := defaultGasLimit
	if logLevel == MediumLogLevel {
		fmt.Println(color.CyanString("usedGasLimit -> defaultGasLimit: "), color.YellowString(strconv.Itoa(int(usedGasLimit)))+"\n")
	}

	if customGasLimit != nil {
		usedGasLimit = customGasLimit.(uint64)

		if logLevel == MediumLogLevel {
			fmt.Println(color.CyanString("usedGasLimit -> customGasLimit: "), color.YellowString(strconv.Itoa(int(usedGasLimit)))+"\n")
		}
	} else {
		if len(data) > 0 {
			usedGasLimit = w.EstimateGas(toAddressString, data)
			if logLevel == MediumLogLevel {
				fmt.Println(color.CyanString("usedGasLimit -> w.EstimateGas: "), color.YellowString(strconv.Itoa(int(usedGasLimit)))+"\n")
			}
		} else {

		}
	}

	toAddress := common.HexToAddress(toAddressString)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: usedGasPrice,
		Gas:      usedGasLimit,
		To:       &toAddress,
		Value:    value,
		Data:     data,
	})

	singedTx, signTxErr := w.SignTx(tx)
	if signTxErr != nil {
		return "", 0, signTxErr
	}

	sendTxErr := w.selectClient().SendTransaction(context.Background(), singedTx)
	if sendTxErr != nil {
		return "", 0, sendTxErr
	}

	if logLevel == HighLogLevel {

		b, e := singedTx.MarshalJSON()
		if e != nil {
			fmt.Println("SendTransaction")
			return "", 0, e
		}

		var result map[string]interface{}
		json.Unmarshal(b, &result)
		s, _ := prettyjson.Marshal(result)

		timestamp := time.Now().Unix()

		fmt.Println(color.GreenString("Raw Transaction Hash: "), color.YellowString(tx.Hash().Hex()))
		fmt.Println(color.CyanString("Transaction Hash: "), color.YellowString(singedTx.Hash().Hex()))
		fmt.Println(color.MagentaString("Timestamp: "), color.YellowString(strconv.Itoa(int(timestamp))))
		fmt.Println(string(s))

		//OpenBrowser("https://testnet.snowtrace.io/tx/" + singedTx.Hash().Hex())
	}

	return singedTx.Hash().Hex(), nonce, nil
}

func (w *Web3Manager) CancelTx(to string, nonce uint64, multiplier int64) (string, error) {

	gasPrice, _ := w.selectClient().SuggestGasPrice(context.Background())

	txId, _, err := w.SignAndSendTransaction(
		to,
		web3util.ToWei(0, 0),
		make([]byte, 0),
		nonce,
		nil,
		big.NewInt(gasPrice.Int64()*multiplier))

	if err != nil {
		return "", err
	}

	return txId, nil
}

func (w *Web3Manager) SendEth(to string, val float64) {
	fromAddress := crypto.PubkeyToAddress(*w.PublicKeyECDSA)
	toAddress := common.HexToAddress(to)
	var data []byte
	nonce, err := w.Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasLimit := uint64(2100000) // in units
	gasPrice, err := w.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	ethValue := ethutils.EtherToWei(big.NewFloat(val))
	tx := types.NewTransaction(nonce, toAddress, ethValue, gasLimit, gasPrice, data)

	chainID, err := w.Client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), w.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}
func (w *Web3Manager) SwitchAccount(plainPrivateKey string) {
	// create privateKey from string key
	privateKey, privateKeyErr := crypto.HexToECDSA(plainPrivateKey)
	errorsutil.HandleError(privateKeyErr)

	// generate public key and address from private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	// generate address from public key
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	w.Address = address
}

func (w *Web3Manager) ConfigureTransactor(value *big.Int, gasPrice *big.Int, gasLimit uint64) {

	if value.String() != "-1" {
		w.Transactor.Value = value
	}

	w.Transactor.GasPrice = gasPrice
	w.Transactor.GasLimit = gasLimit
	w.Transactor.Nonce = w.PendingNonce()
	w.Transactor.Context = context.Background()
}

func (w *Web3Manager) Balance() *big.Int {
	// get current balance
	balance, balanceErr := w.Client.BalanceAt(context.Background(), w.Address, nil)
	errorsutil.HandleError(balanceErr)
	return balance
}

func (w *Web3Manager) PendingNonce() *big.Int {
	// calculate next nonce
	nonce, nonceErr := w.Client.PendingNonceAt(context.Background(), w.Address)
	errorsutil.HandleError(nonceErr)
	return big.NewInt(int64(nonce))
}
