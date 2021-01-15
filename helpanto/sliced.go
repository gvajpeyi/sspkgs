package helpanto

type StringSlice []string


func (s *StringSlice) In(a string) bool {
	ss := *s
	
	for i := 0; i < len(ss); i++ {
		if a == ss[i] {
			return true
		}
	}
	return false
	
}


