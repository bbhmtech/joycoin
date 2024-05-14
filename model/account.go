package model

import (
	"bytes"
	"crypto/sha256"
	"time"

	"github.com/google/uuid"
)

const myPhrase = `So big smiles, we're gonna get in front of this`

type Account struct {
	ID                uint `gorm:"primarykey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Role              string
	Activated         bool
	CachedCentBalance int64
	PasscodeHash      []byte
	DeviceBindingKey  string `gorm:"index"`
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

func (a *Account) NewDeviceBindingKey() {
	a.DeviceBindingKey = uuid.NewString()
}
