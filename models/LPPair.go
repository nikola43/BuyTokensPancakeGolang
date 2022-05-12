package models

import (
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

type LPPair struct {
	gorm.Model
	LPAddress    common.Address
	LPPairA      common.Address
	LPPairB      common.Address
	HasLiquidity bool
}
