package models

import "gorm.io/gorm"


type LPPair struct {
	gorm.Model
	TxHash       string
	LPAddress    string
	LPPairA      string
	LPPairB      string
	HasLiquidity bool
}