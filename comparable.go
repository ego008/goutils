package goutils

// ComparableContains Check if a slice contain an element
// Usage:
//	var stringSlice = []string{"item 1", "item 2"}
//	var intSlice = []int{1, 2, 3}
//
//	fmt.Println(ComparableContains(stringSlice, "item 2")) // true
//	fmt.Println(ComparableContains(intSlice, 4)) // false
func ComparableContains[T comparable](slice []T, element T) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}

	return false
}

// ComparableFilter Filter a slice to return only the elements that match a condition
// Usage:
//	var intSlice = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
//	var oddSlice = ComparableFilter(intSlice, func(element int) bool {
//    return element % 2 != 0
//	})
//	fmt.Println(oddSlice) // []int{1, 3, 5, 7, 9}
func ComparableFilter[T comparable](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, e := range slice {
		if predicate(e) {
			result = append(result, e)
		}
	}

	return result
}

// ComparableMap Transform a slice into a difference slice type
// Usage:
//	var intSlice = []int{1, 2, 3}
//	var strSlice = ComparableMap(intSlice, func(element int) string {
//    return strconv.Itoa(element)
//	})
//	fmt.Println(strSlice) // []string{"1", "2", "3"}
func ComparableMap[T comparable, R comparable](slice []T, mapper func(T) R) []R {
	var result []R
	for _, e := range slice {
		result = append(result, mapper(e))
	}

	return result
}

// ComparableRepeat Create a slice with a repeated values
// Usage:
//	var intSlice = ComparableRepeat(1, 3)
//	fmt.Println(intSlice) // []int{1, 1, 1}
func ComparableRepeat[T comparable](input T, times int) []T {
	result := make([]T, times)
	for i := 0; i < times; i++ {
		result[i] = input
	}

	return result
}

// ComparableOverlap Get the overlap elements between two slice
// Usage:
//	var slice1 = []int{1, 2, 3}
//	var slice2 = []int{2, 3, 4}
//	var result = ComparableOverlap(slice1, slice2)
//	fmt.Println(result) // []int{2, 3}
func ComparableOverlap[T comparable](slice1 []T, slice2 []T) []T {
	result := make([]T, 0)
	for _, e1 := range slice1 {
		if ComparableContains(slice2, e1) {
			result = append(result, e1)
		}
	}

	return result
}

// ComparableAppendIfNotExists Append an element to a slice only if it isn’t existed
// Usage:
//	var intSlice = []int{1, 2, 3}
//	var intSlice1 = ComparableAppendIfNotExists(intSlice, 3)
//	var intSlice2 = ComparableAppendIfNotExists(intSlice, 4)
//	fmt.Println(intSlice1) // []int{1, 2, 3}
//	fmt.Println(intSlice2) // []int{1, 2, 3, 4}
func ComparableAppendIfNotExists[T comparable](slice []T, newItem T) []T {
	for _, s := range slice {
		if s == newItem {
			return slice
		}
	}

	return append(slice, newItem)
}

// AnyAppendIfNotExists Append an element to a slice only if it isn’t existed
// Usage:
//	type Thing struct {
//   ID int
//   Name string
//	}
//	var things = []Thing{
//	{ID: 1, Name: "one"},
//	{ID: 2, Name: "two"},
//	}
//	var thing2 = Thing{ID: 2, Name: "Thing 2"}
//	var thingSlice1 = AnyAppendIfNotExists(things, thing2, func (element Thing) bool {
//   return element.ID == thing2.ID
//	})
//	fmt.Println(len(thingSlice1)) // 2
func AnyAppendIfNotExists[T any](slice []T, newItem T, checkExists func(T) bool) []T {
	for k := range slice {
		if checkExists(slice[k]) {
			return slice
		}
	}

	return append(slice, newItem)
}
