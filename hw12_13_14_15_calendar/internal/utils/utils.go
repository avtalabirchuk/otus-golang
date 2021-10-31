package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateUID() (res int64, err error) {
	n, err := rand.Int(rand.Reader, big.NewInt(100000))
	if err != nil {
		return
	}
	return n.Int64(), nil
}
