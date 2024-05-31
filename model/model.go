package model

import (
	"log/slog"

	"gorm.io/gorm"
)

func createRootAccount(db *gorm.DB) {

	acc := Account{
		Role:              "operator",
		Activated:         false,
		CachedCentBalance: 0,
	}

	if err := db.Save(&acc).Error; err != nil {
		panic(err)
	}

	slog.Info("create jumpber for root", "id", acc.ID, "role", acc.Role)
	j, err := CreateJumpberFromAccount(db, &acc)
	if err != nil {
		panic(err)
	}

	eID, err := j.EncodeID()
	if err != nil {
		panic(err)
	}
	slog.Info("jumper created", "id", j.ID, "encoded", eID, "err", err)
}

func AutoMigration(db *gorm.DB) {
	slog.Info("execute: db.AutoMigrate")
	err := db.AutoMigrate(&Jumper{}, &Account{}, &ShortenLink{}, &Transaction{}, &QuickAction{})
	slog.Info("migration complete", "err", err)

	var nOperator int64 = 0
	err = db.Model(&Account{}).Where("role = 'operator'").Count(&nOperator).Error
	slog.Debug("counting role=operator", "number", nOperator, "err", err)
	if err != nil {
		panic(err)
	}

	if nOperator == 0 {
		slog.Info("creating operator account", "nOperator", nOperator)
		createRootAccount(db)
	}
}
