package models

import (
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

type EventsCatched struct {
	gorm.Model
	TxHash       string
	TokenAddress common.Address
	LPPairs      []LPPair
}
