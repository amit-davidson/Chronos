package utils

type IntSet struct {
	set map[int]int
}


func (s *IntSet) Add(num int) bool {
	if !s.Exist(num) {
		s.set[num] = 1
		return true
	}
	return false
}

func (s *IntSet) Exist(num int) bool {
	_, exist := s.set[num]
	return exist
}
