package utils

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"time"
)

func NewUlidString() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}
