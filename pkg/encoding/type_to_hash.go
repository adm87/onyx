package encoding

import (
	"crypto/sha256"
	"encoding/binary"
	"reflect"
)

func TypeID[T any]() uint64 {
	sum := sha256.Sum256([]byte(reflect.TypeFor[T]().String()))
	return binary.BigEndian.Uint64(sum[:8])
}
