package model

import "time"

type QuickAction struct {
	DeviceBindingKey string `gorm:"primarykey"`
	ValidBefore      time.Time
	Temporary        bool
	CachedAccountID  uint
	CachedAccount    Account `gorm:"foreignKey:CachedAccountID"`
	Action           string
	Int64Value1      int64
	UintValue1       uint
	StringValue1     string
}

func (qa *QuickAction) IsValid() bool {
	return time.Now().Before(qa.ValidBefore)
}
