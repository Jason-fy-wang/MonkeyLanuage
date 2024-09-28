package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
)

type HashTable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value int64
}

func (b *Boolean) HashKey() HashKey {
	var val int64
	if b.Value {
		val = 1
	} else {
		val = 0
	}
	return HashKey{Type: b.Type(), Value: val}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: int64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()

	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: int64(h.Sum64())}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType {
	return HASH_OBJ
}

func(h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string[]

	for _,  pair:= range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s:%s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs,", "))
	out.WriteString("}")
	return out.String()
}
