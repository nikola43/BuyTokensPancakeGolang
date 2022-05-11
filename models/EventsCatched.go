package models

import "gorm.io/gorm"


type EventsCatched struct {
	gorm.Model
	TxHash       string
	TokenAddress string
	LPAddress    string
	LPPairA      string
	LPPairB      string
	Timestamp    string
	HasLiquidity bool
}