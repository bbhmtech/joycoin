package model

import (
	"bytes"
	"crypto/sha256"
	"time"
)

const myPhrase = `So big smiles, we're gonna get in front of this`

type Account struct {
	ID                uint      `json:"id" gorm:"primarykey"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Nickname          string    `json:"nickname"`
	Role              string    `json:"role"`
	Activated         bool      `json:"activated"`
	CachedCentBalance int64     `json:"cached_cent_balance"`
	PasscodeHash      []byte    `json:"-"`
	DeviceBindingKey  string    `json:"-"  gorm:"uniqueIndex"`
}

func (a *Account) ChangePasscode(new string) {
	h := sha256.New()
	h.Write([]byte(new))
	h.Write([]byte(myPhrase))
	a.PasscodeHash = h.Sum(nil)
}

func (a *Account) VerifyPasscode(code string) bool {
	h := sha256.New()
	h.Write([]byte(code))
	h.Write([]byte(myPhrase))
	return bytes.Equal(a.PasscodeHash, h.Sum(nil))
}

func (a *Account) IsMerchant() bool {
	return a.Role == "merchant"
}

func (a *Account) IsNormal() bool {
	return a.Role == "normal"
}

func (a *Account) IsOperator() bool {
	return a.Role == "operator"
}
