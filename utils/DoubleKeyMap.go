package utils

import (
	"go/token"
	"hash/fnv"
	"strconv"
)

type DoubleKeyMap map[uint32]struct{}

func NewDoubleKeyMap() DoubleKeyMap {
	return make(DoubleKeyMap, 0)
}

func getHash(num int) uint32 {
	h := fnv.New32a()
	h.Write([]byte(strconv.Itoa(num)))
	return h.Sum32()
}

func calcHash(posA token.Pos, posB token.Pos) uint32 {
	hashA := getHash(int(posA))
	hashB := getHash(int(posB))
	return hashA+hashB
}


func (m DoubleKeyMap) Add(posA token.Pos, posB token.Pos) {
	key := calcHash(posA, posB)
	m[key] = struct{}{} // Hashes are added to make the dict commutative

}
func (m DoubleKeyMap) IsExist(posA token.Pos, posB token.Pos) bool {
	key := calcHash(posA, posB)
	_, ok := m[key]
	return ok
}
