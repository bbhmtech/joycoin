package model

import (
	"github.com/google/uuid"
	"github.com/mr-tron/base58"
	"gorm.io/gorm"
)

type Jumper struct {
	ID       string
	Hint     string
	TargetID uint
}

func GetJumperFromEncodedID(db *gorm.DB, key string) (*Jumper, error) {
	b, err := base58.Decode(key)
	if err != nil {
		return nil, err
	}

	id, err := uuid.FromBytes(b)
	if err != nil {
		return nil, err
	}

	j := Jumper{ID: id.String()}
	err = db.Take(&j).Error
	return &j, err
}

func CreateJumpberFromAccount(db *gorm.DB, acc *Account) (*Jumper, error) {
	j := Jumper{ID: uuid.NewString(), Hint: "NTAG|Account", TargetID: acc.ID}
	err := db.Save(&j).Error
	return &j, err
}

func (j *Jumper) EncodeID() (string, error) {
	u, err := uuid.Parse(j.ID)
	if err != nil {
		return "", err
	}

	b, err := u.MarshalBinary()
	if err != nil {
		return "", err
	}
	return base58.Encode(b), nil
}
