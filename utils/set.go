package utils

type IntSet struct {
	set map[int]int
}

// NewStringSet - creates a new strings set
func NewIntSet() *IntSet {
	ss := &IntSet{set: make(map[int]int)}
	return ss
}

func (ss *IntSet) Add(num int) bool {
	if !ss.Exist(num) {
		ss.set[num] = 1
		return true
	}
	return false
}

func (ss *IntSet) Exist(num int) bool {
	_, exist := ss.set[num]
	return exist
}
