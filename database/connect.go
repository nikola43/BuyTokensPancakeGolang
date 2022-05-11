package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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

func InitDatabase() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&EventsCatched{})

	// Read
	var product Product
	db.First(&product, 1)                 // find product with integer primary key
	db.First(&product, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	db.Model(&product).Update("Price", 200)
	// Update - update multiple fields
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	db.Delete(&product, 1)
}

func InsertNewEvent(newEvent *EventsCatched) bool {
	res := db.Create(&EventsCatched{TxHash: newEvent.TxHash, TokenAddress: newEvent.TokenAddress, LPAddress: newEvent.LPAddress, LPPairA: newEvent.LPPairA, LPPairB: newEvent.LPPairB, Timestamp: newEvent.Timestamp, HasLiquidity: false})

	return res
}

func UpdateLiquidity(txHash string) bool {
	var event EventsCatched
	db.First(&event, "TxHash = ?", txHash)

	return res
}
