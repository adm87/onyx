package engine

import (
	"hash/fnv"
	"reflect"
)

func TypeHash[T any]() uint64 {
	h := fnv.New64a()
	t := reflect.TypeFor[T]()
	h.Write([]byte(t.PkgPath() + "." + t.Name()))
	return h.Sum64()
}
