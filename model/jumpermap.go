package model

import (
	"github.com/google/uuid"
	"github.com/mr-tron/base58"
)

type JumperMap struct {
	ID        string
	Hint      string
	TargetURL string
}

func JumperMapFromEncodedID(key string) (*JumperMap, error) {
	b, err := base58.Decode(key)
	if err != nil {
		return nil, err
	}

	id, err := uuid.FromBytes(b)
	if err != nil {
		return nil, err
	}

	return &JumperMap{
		ID:        id.String(),
		TargetURL: "demo",
	}, nil
}

func (j *JumperMap) EncodeID() (string, error) {
	u, err := uuid.Parse(j.ID)
	if err != nil {
		return "", err
	}

	b, err := u.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return base58.Encode(b), nil
}
