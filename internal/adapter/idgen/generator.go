package idgen

import (
	"crypto/rand"
	"strings"
)

type GenID struct {
	len int
}

func NewGenID(l int) *GenID {
	return &GenID{
		len: l,
	}
}

func (g *GenID) BookRef() (string, error) {
	const alnum = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	buf := make([]byte, g.len)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	var sb strings.Builder

	for _, b := range buf {
		symbol := alnum[int(b)%len(alnum)]
		sb.WriteByte(symbol)
	}

	return sb.String(), nil
}
