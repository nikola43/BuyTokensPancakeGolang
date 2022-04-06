package ethbasedclient

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nikola43/BuyTokensPancakeGolang/errorsutil"
	"github.com/nikola43/BuyTokensPancakeGolang/ethutils"
	"log"
	"math/big"
)

type EthBasedClient struct {
	Client         *ethclient.Client
	PrivateKey     *ecdsa.PrivateKey
	PublicKeyECDSA *ecdsa.PublicKey
	Address        common.Address
	ChainID        *big.Int
	Transactor     *bind.TransactOpts
	Nonce          *big.Int
}

func New(rawurl, plainPrivateKey string) EthBasedClient {
	// connect with bsc
	client, err := ethclient.Dial(rawurl)
	errorsutil.HandleError(err)

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

	// get chain id
	chainID, chainIDErr := client.ChainID(context.Background())
	errorsutil.HandleError(chainIDErr)

	// generate transactor for transactions management
	transactor, transactOptsErr := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	errorsutil.HandleError(transactOptsErr)

	ethBasedClientTemp := EthBasedClient{
		Client:         client,
		PrivateKey:     privateKey,
		PublicKeyECDSA: publicKeyECDSA,
		Address:        address,
		ChainID:        chainID,
		Transactor:     transactor,
	}

	return ethBasedClientTemp
}

func (ethBasedClient *EthBasedClient) SendEth(to string, val float64) {
	fromAddress := crypto.PubkeyToAddress(*ethBasedClient.PublicKeyECDSA)
	toAddress := common.HexToAddress(to)
	var data []byte
	nonce, err := ethBasedClient.Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasLimit := uint64(2100000) // in units
	gasPrice, err := ethBasedClient.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	ethValue := ethutils.EtherToWei(big.NewFloat(val))
	tx := types.NewTransaction(nonce, toAddress, ethValue, gasLimit, gasPrice, data)

	chainID, err := ethBasedClient.Client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), ethBasedClient.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = ethBasedClient.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}
func (ethBasedClient *EthBasedClient) SwitchAccount(plainPrivateKey string) {
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
	ethBasedClient.Address = address
}

func (ethBasedClient *EthBasedClient) ConfigureTransactor(value *big.Int, gasPrice *big.Int, gasLimit uint64) {

	if value.String() != "-1" {
		ethBasedClient.Transactor.Value = value
	}

	ethBasedClient.Transactor.GasPrice = gasPrice
	ethBasedClient.Transactor.GasLimit = gasLimit
	ethBasedClient.Transactor.Nonce = ethBasedClient.PendingNonce()
	ethBasedClient.Transactor.Context = context.Background()
}

func (ethBasedClient *EthBasedClient) Balance() *big.Int {
	// get current balance
	balance, balanceErr := ethBasedClient.Client.BalanceAt(context.Background(), ethBasedClient.Address, nil)
	errorsutil.HandleError(balanceErr)
	return balance
}

func (ethBasedClient *EthBasedClient) PendingNonce() *big.Int {
	// calculate next nonce
	nonce, nonceErr := ethBasedClient.Client.PendingNonceAt(context.Background(), ethBasedClient.Address)
	errorsutil.HandleError(nonceErr)
	return big.NewInt(int64(nonce))
}
