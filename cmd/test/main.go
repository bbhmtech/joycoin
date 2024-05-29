package main

import (
	"fmt"

	"github.com/bbhmtech/joycoin/model"
)

func main() {
	v := model.Jumper{
		ID:       "54422ac5-8045-4dbe-99b5-180b2340b962",
		Hint:     "NTAG|Account",
		TargetID: 3,
	}
	fmt.Println(v.EncodeID())
}
