package model

import (
	"time"
)

type QuickAction struct {
	DeviceBindingKey string    `json:"-" gorm:"primarykey"`
	ValidBefore      time.Time `json:"valid_before"`
	Temporary        bool      `json:"temporary"`
	CachedAccountID  uint      `json:"cached_account_id"`
	CachedAccount    Account   `json:"-" gorm:"foreignKey:CachedAccountID"`
	Action           string    `json:"action"`
	Int64Value1      int64     `json:"int64_value_1"`
	UintValue1       uint      `json:"uint_value_1"`
	StringValue1     string    `json:"string_value_1"`
}

func (qa *QuickAction) IsValid() bool {
	return time.Now().Before(qa.ValidBefore) && qa.CachedAccount.DeviceBindingKey == qa.DeviceBindingKey
}
