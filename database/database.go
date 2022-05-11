package controllers

import (
	"gorm.io/gorm"
)

var GormDB *gorm.DB

func Migrate() {
	// DROP
	//GormDB.Migrator().DropTable(&models.User{})
	//GormDB.Migrator().DropTable(&models.NodeDB{})
	//GormDB.Migrator().DropTable(&models.Gift{})
	//GormDB.Migrator().DropTable(&models.GiftCard{})

	// CREATE
	//GormDB.AutoMigrate(&models.User{})
	//GormDB.AutoMigrate(&models.NodeDB{})
	//GormDB.AutoMigrate(&models.Gift{})
	//GormDB.AutoMigrate(&models.GiftCard{})
}
