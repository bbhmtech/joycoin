package main

import (
	"github.com/bbhmtech/joycoin"
	"github.com/bbhmtech/joycoin/model"
)

func main() {

	cfg := joycoin.LoadConfig("config.json")
	db := cfg.InitializeDatabase()
	model.AutoMigration(db)

	acc := model.Account{
		ID:                1,
		Role:              "operator",
		Activated:         false,
		CachedCentBalance: 0,
	}

	if err := db.Save(&acc).Error; err != nil {
		panic(err)
	}
}
