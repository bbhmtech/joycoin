package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID           uint   `gorm:"primarykey"`
	ReferenceTag string `gorm:"uniqueIndex"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	InitiatorAccountID uint
	InitiatorAccount   Account `json:"-" gorm:"foreignKey:InitiatorAccountID"`
	FromAccountID      uint
	FromAccount        Account `json:"-" gorm:"foreignKey:FromAccountID"`
	ToAccountID        uint
	ToAccount          Account `json:"-" gorm:"foreignKey:ToAccountID"`
	Message            string
	CentAmount         int64
	// from: -amount, to: +amount
}

func (t *Transaction) PreFlightCheck(db *gorm.DB) error {
	err := db.Take(&t.InitiatorAccount, t.InitiatorAccountID).Error
	if err != nil {
		return fmt.Errorf("无法找到 Initiator: %w", err)
	}

	if !t.InitiatorAccount.IsOperator() && t.InitiatorAccount.ID != t.FromAccount.ID {
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

	if t.FromAccount.IsNormal() && t.CentAmount > 0 && t.FromAccount.CachedCentBalance >= t.CentAmount {
		return errors.New("支出方账户余额不足")
	}

	return nil
}
