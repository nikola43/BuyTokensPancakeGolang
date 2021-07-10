package ethbasedclient

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nikola43/buy_pancake/utils/errorsutil"
	"github.com/nikola43/buy_pancake/utils/ethutil"
	"math/big"
)

type EthBasedClient struct {
	Client     *ethclient.Client
	PrivateKey *ecdsa.PrivateKey
	Address    common.Address
	ChainID    *big.Int
	Transactor *bind.TransactOpts
	Nonce      *big.Int
}

func New(rawUrl, plainPrivateKey string) EthBasedClient {
	client, err := ethclient.Dial(rawUrl)
	errorsutil.HandleError(err)

	privateKey := ethutil.GenerateEcdsaPrivateKey(plainPrivateKey)
	ethBasedClientTemp := EthBasedClient{
		Client:     client,
		PrivateKey: privateKey,
		Address:    ethutil.GenerateAddress(privateKey),
		ChainID:    ethutil.GetChainID(client),
		Transactor: ethutil.GenerateTransactor(client, privateKey),
	}

	return ethBasedClientTemp
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
	balance, balanceErr := ethBasedClient.Client.BalanceAt(context.Background(), ethBasedClient.Address, nil)
	errorsutil.HandleError(balanceErr)
	return balance
}

func (ethBasedClient *EthBasedClient) PendingNonce() *big.Int {
	nonce, nonceErr := ethBasedClient.Client.PendingNonceAt(context.Background(), ethBasedClient.Address)
	errorsutil.HandleError(nonceErr)
	return big.NewInt(int64(nonce))
}
