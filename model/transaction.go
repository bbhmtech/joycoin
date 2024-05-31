package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID           uint      `gorm:"primarykey"`
	ReferenceTag string    `gorm:"uniqueIndex"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	InitiatorAccountID uint    `json:"initiator_account_id"`
	InitiatorAccount   Account `json:"-" gorm:"foreignKey:InitiatorAccountID"`
	FromAccountID      uint    `json:"from_account_id"`
	FromAccount        Account `json:"-" gorm:"foreignKey:FromAccountID"`
	ToAccountID        uint    `json:"to_account_id"`
	ToAccount          Account `json:"-" gorm:"foreignKey:ToAccountID"`
	Message            string  `json:"message"`
	CentAmount         int64   `json:"cent_amount"`
	// from: -amount, to: +amount
}

func (t *Transaction) PreFlightCheck(db *gorm.DB) error {
	if t.FromAccountID == t.ToAccountID {
		return errors.New("无法进行同账户交易")
	}

	err := db.Take(&t.InitiatorAccount, t.InitiatorAccountID).Error
	if err != nil {
		return fmt.Errorf("无法找到 Initiator: %w", err)
	}

	if t.InitiatorAccount.IsNormal() && t.InitiatorAccountID != t.FromAccountID {
		return errors.New("未授权的访问")
	}

	err = db.Take(&t.FromAccount, t.FromAccountID).Error
	if err != nil {
		return fmt.Errorf("无法找到支出方账户: %w", err)
	}

	if !t.FromAccount.Activated {
		return errors.New("支出方账户未激活")
	}

	err = db.Take(&t.ToAccount, t.ToAccountID).Error
	if err != nil {
		return fmt.Errorf("无法找到接收方账户: %w", err)
	}

	if !t.ToAccount.Activated {
		return errors.New("接收方账户未激活")
	}

	// I know here racing condition exists, but I wanna to prioritize success in transactions first.
	// If negative balance occurs, subsequent transaction will fail.
	if t.FromAccount.IsNormal() && t.CentAmount > 0 && t.FromAccount.CachedCentBalance < t.CentAmount {
		return errors.New("支出方账户余额不足")
	}

	return nil
}

func (t *Transaction) Insert(db *gorm.DB) error {
	if t.ID != 0 {
		return errors.New("already inserted")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&t).Error; err != nil {
			return err
		}

		if err := tx.Model(&t.FromAccount).Update("cached_cent_balance", gorm.Expr("cached_cent_balance - ?", t.CentAmount)).Error; err != nil {
			return err
		}

		if err := tx.Model(&t.ToAccount).Update("cached_cent_balance", gorm.Expr("cached_cent_balance + ?", t.CentAmount)).Error; err != nil {
			return err
		}

		return nil
	})
}
