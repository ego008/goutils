package lst

import "sort"

type (
	StrLst []string
	IntLst []int
	Mp     map[int]struct{}
)

func (l StrLst) Has(x string) bool {
	for _, i := range l {
		if i == x {
			return true
		}
	}
	return false
}

func (l StrLst) Sort() {
	var s []string
	for _, i := range l {
		s = append(s, i)
	}

	sort.Strings(s) // 升序

	for i, v := range s {
		l[i] = v
	}
}

func (l IntLst) Has(x int) bool {
	for _, i := range l {
		if i == x {
			return true
		}
	}
	return false
}

func (l IntLst) Sort() {
	var s []int
	for _, i := range l {
		s = append(s, int(i))
	}

	sort.Ints(s) // 升序
	// sort.Sort(sort.Reverse(sort.IntSlice(s))) // 降序

	for i, v := range s {
		l[i] = v
	}
}

func (l IntLst) ReverseSort() {
	var s []int
	for _, i := range l {
		s = append(s, int(i))
	}

	sort.Sort(sort.Reverse(sort.IntSlice(s))) // 降序

	for i, v := range s {
		l[i] = v
	}
}

func (m Mp) Has(x int) (ok bool) {
	_, ok = m[x]
	return
}
