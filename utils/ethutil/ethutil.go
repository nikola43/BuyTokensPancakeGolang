package ethutil

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/nikola43/buy_pancake/utils/errorsutil"
	"github.com/shopspring/decimal"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func GweiToEther(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(wei), big.NewFloat(params.GWei))
}

func GweiToWei(wei *big.Int) *big.Int {
	return EtherToWei(GweiToEther(wei))
}

// Wei ->
func WeiToGwei(wei *big.Int) *big.Int {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	v := f.Quo(fWei.SetInt(wei), big.NewFloat(params.GWei))
	i, _ := new(big.Int).SetString(v.String(), 10)
	return i
}

func WeiToEther(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(wei), big.NewFloat(params.Ether))
}

// ETH -> Wei
func EtherToWei(eth *big.Float) *big.Int {
	s, err := strconv.ParseFloat(eth.String(), 64)
	errorsutil.HandleError(err)
	return ToWei(s, 18)
}

func EtherToGwei(eth *big.Float) *big.Int {
	truncInt, _ := eth.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.GWei))
	fracStr := strings.Split(fmt.Sprintf("%.9f", eth), ".")[1]
	fracStr += strings.Repeat("0", 9-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}

func ToWei(iamount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}

// IsValidAddress validate hex address
func IsValidAddress(iaddress interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := iaddress.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

// IsZeroAddress validate if it's a 0 address
func IsZeroAddress(iaddress interface{}) bool {
	var address common.Address
	switch v := iaddress.(type) {
	case string:
		address = common.HexToAddress(v)
	case common.Address:
		address = v
	default:
		return false
	}

	zeroAddressBytes := common.FromHex("0x0000000000000000000000000000000000000000")
	addressBytes := address.Bytes()
	return reflect.DeepEqual(addressBytes, zeroAddressBytes)
}

// ToDecimal wei to decimals
func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

// CalcGasCost calculate gas cost given gas limit (units) and gas price (wei)
func CalcGasCost(gasLimit uint64, gasPrice *big.Int) *big.Int {
	return big.NewInt(0).Mul(big.NewInt(int64(gasLimit)), gasPrice)
}

func GeneratePath(tokenAContractPlainAddress string, tokenBContractPlainAddress string) []common.Address {
	tokenAContractAddress := common.HexToAddress(tokenAContractPlainAddress)
	tokenBContractAddress := common.HexToAddress(tokenBContractPlainAddress)

	path := make([]common.Address, 0)
	path = append(path, tokenAContractAddress)
	path = append(path, tokenBContractAddress)

	return path
}

func GetChainID(client *ethclient.Client) *big.Int {
	// get chain id
	chainID, err := client.ChainID(context.Background())
	errorsutil.HandleError(err)
	return chainID
}

func CalculatePercent(value *big.Int, percent int64) *big.Int {
	return big.NewInt(0).Div(big.NewInt(0).Mul(value, big.NewInt(percent)), big.NewInt(100))
}

func CancelTransaction(client *ethclient.Client, transaction *types.Transaction, privateKey *ecdsa.PrivateKey) (*types.Transaction, error) {
	var txData []byte

	newGasPrice := big.NewInt(0).Add(transaction.GasPrice(), CalculatePercent(transaction.GasPrice(), 10))
	tx := types.NewTransaction(
		transaction.Nonce(),
		GenerateAddress(privateKey),
		big.NewInt(0),
		transaction.Gas(),
		newGasPrice,
		txData)

	signedTx := SendTransaction(client, tx, privateKey)

	return signedTx, nil
}

func SendTransaction(client *ethclient.Client, tx *types.Transaction, privateKey *ecdsa.PrivateKey) *types.Transaction {
	signedTx := SignTransaction(client, tx, privateKey)

	err := client.SendTransaction(context.Background(), signedTx)
	errorsutil.HandleError(err)

	return signedTx
}

func GenerateTransactor(client *ethclient.Client, privateKey *ecdsa.PrivateKey) *bind.TransactOpts {
	chainID := GetChainID(client)

	transactor, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	errorsutil.HandleError(err)

	return transactor
}

func GenerateEcdsaPrivateKey(plainPrivateKey string) *ecdsa.PrivateKey {
	privateKey, privateKeyErr := crypto.HexToECDSA(plainPrivateKey)
	errorsutil.HandleError(privateKeyErr)

	return privateKey
}

func GenerateAddress(privateKey *ecdsa.PrivateKey) common.Address {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		errorsutil.HandleError(errors.New("error casting public key to ECDSA"))
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

func SignTransaction(client *ethclient.Client, tx *types.Transaction, privateKey *ecdsa.PrivateKey) *types.Transaction {
	chainID := GetChainID(client)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	errorsutil.HandleError(err)

	return signedTx
}
