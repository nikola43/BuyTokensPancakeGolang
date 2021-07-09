// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package PancakeLibrary

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// PancakeLibraryABI is the input ABI used to generate the binding from.
const PancakeLibraryABI = "[]"

// PancakeLibraryBin is the compiled bytecode used for deploying new contracts.
var PancakeLibraryBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea264697066735822122079a83e13bf586287183b58e5479ec5b0557e0df16606453cffb7c7fd65ae4aa464736f6c63430006060033"

// DeployPancakeLibrary deploys a new Ethereum contract, binding an instance of PancakeLibrary to it.
func DeployPancakeLibrary(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *PancakeLibrary, error) {
	parsed, err := abi.JSON(strings.NewReader(PancakeLibraryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(PancakeLibraryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PancakeLibrary{PancakeLibraryCaller: PancakeLibraryCaller{contract: contract}, PancakeLibraryTransactor: PancakeLibraryTransactor{contract: contract}, PancakeLibraryFilterer: PancakeLibraryFilterer{contract: contract}}, nil
}

// PancakeLibrary is an auto generated Go binding around an Ethereum contract.
type PancakeLibrary struct {
	PancakeLibraryCaller     // Read-only binding to the contract
	PancakeLibraryTransactor // Write-only binding to the contract
	PancakeLibraryFilterer   // Log filterer for contract events
}

// PancakeLibraryCaller is an auto generated read-only Go binding around an Ethereum contract.
type PancakeLibraryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PancakeLibraryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PancakeLibraryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PancakeLibraryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PancakeLibraryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PancakeLibrarySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PancakeLibrarySession struct {
	Contract     *PancakeLibrary   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PancakeLibraryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PancakeLibraryCallerSession struct {
	Contract *PancakeLibraryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// PancakeLibraryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PancakeLibraryTransactorSession struct {
	Contract     *PancakeLibraryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// PancakeLibraryRaw is an auto generated low-level Go binding around an Ethereum contract.
type PancakeLibraryRaw struct {
	Contract *PancakeLibrary // Generic contract binding to access the raw methods on
}

// PancakeLibraryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PancakeLibraryCallerRaw struct {
	Contract *PancakeLibraryCaller // Generic read-only contract binding to access the raw methods on
}

// PancakeLibraryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PancakeLibraryTransactorRaw struct {
	Contract *PancakeLibraryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPancakeLibrary creates a new instance of PancakeLibrary, bound to a specific deployed contract.
func NewPancakeLibrary(address common.Address, backend bind.ContractBackend) (*PancakeLibrary, error) {
	contract, err := bindPancakeLibrary(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PancakeLibrary{PancakeLibraryCaller: PancakeLibraryCaller{contract: contract}, PancakeLibraryTransactor: PancakeLibraryTransactor{contract: contract}, PancakeLibraryFilterer: PancakeLibraryFilterer{contract: contract}}, nil
}

// NewPancakeLibraryCaller creates a new read-only instance of PancakeLibrary, bound to a specific deployed contract.
func NewPancakeLibraryCaller(address common.Address, caller bind.ContractCaller) (*PancakeLibraryCaller, error) {
	contract, err := bindPancakeLibrary(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PancakeLibraryCaller{contract: contract}, nil
}

// NewPancakeLibraryTransactor creates a new write-only instance of PancakeLibrary, bound to a specific deployed contract.
func NewPancakeLibraryTransactor(address common.Address, transactor bind.ContractTransactor) (*PancakeLibraryTransactor, error) {
	contract, err := bindPancakeLibrary(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PancakeLibraryTransactor{contract: contract}, nil
}

// NewPancakeLibraryFilterer creates a new log filterer instance of PancakeLibrary, bound to a specific deployed contract.
func NewPancakeLibraryFilterer(address common.Address, filterer bind.ContractFilterer) (*PancakeLibraryFilterer, error) {
	contract, err := bindPancakeLibrary(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PancakeLibraryFilterer{contract: contract}, nil
}

// bindPancakeLibrary binds a generic wrapper to an already deployed contract.
func bindPancakeLibrary(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PancakeLibraryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PancakeLibrary *PancakeLibraryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PancakeLibrary.Contract.PancakeLibraryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PancakeLibrary *PancakeLibraryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PancakeLibrary.Contract.PancakeLibraryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PancakeLibrary *PancakeLibraryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PancakeLibrary.Contract.PancakeLibraryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PancakeLibrary *PancakeLibraryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PancakeLibrary.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PancakeLibrary *PancakeLibraryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PancakeLibrary.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PancakeLibrary *PancakeLibraryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PancakeLibrary.Contract.contract.Transact(opts, method, params...)
}
