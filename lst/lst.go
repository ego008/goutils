package lst

import "sort"

type (
	StrLst []string
	IntLst []int
	Mp     map[int]struct{}
)

func (l StrLst) Has(x string) bool {
	for i := 0; i < len(l); i++ {
		if l[i] == x {
			return true
		}
	}
	return false
}

func (l StrLst) Sort() {
	var s []string
	for i := 0; i < len(l); i++ {
		s = append(s, l[i])
	}

	sort.Strings(s) // 升序

	for i := 0; i < len(s); i++ {
		l[i] = s[i]
	}
}

func (l IntLst) Has(x int) bool {
	for i := 0; i < len(l); i++ {
		if l[i] == x {
			return true
		}
	}
	return false
}

func (l IntLst) Sort() {
	var s []int
	for i := 0; i < len(l); i++ {
		s = append(s, l[i])
	}

	sort.Ints(s) // 升序
	// sort.Sort(sort.Reverse(sort.IntSlice(s))) // 降序

	for i := 0; i < len(s); i++ {
		l[i] = s[i]
	}
}

func (l IntLst) ReverseSort() {
	var s []int
	for i := 0; i < len(l); i++ {
		s = append(s, l[i])
	}

	sort.Sort(sort.Reverse(sort.IntSlice(s))) // 降序

	for i := 0; i < len(s); i++ {
		l[i] = s[i]
	}
}

func (m Mp) Has(x int) (ok bool) {
	_, ok = m[x]
	return
}
