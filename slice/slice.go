package slice

import "reflect"

func ToAnyList[T any](input []T) []any {
	list := make([]any, len(input))
	for i, v := range input {
		list[i] = v
	}
	return list
}

func Difference[T comparable](a []T, b []T) []T {
	set := make([]T, 0)
	for _, v := range a {
		if !Contains(b, v) {
			set = append(set, v)
		}
	}
	return set
}

// IntersectSlices has complexity: O(n^2)
func IntersectSlices[T comparable](a []T, b []T) []T {
	set := make([]T, 0)
	for _, v := range a {
		if Contains(b, v) {
			set = append(set, v)
		}
	}
	return set
}

func DeleteFrom[T comparable](collection []T, el T) []T {
	idx := FindIndex(collection, el)
	if idx > -1 {
		return append(collection[:idx], collection[idx+1:]...)
	}
	return collection
}

func FindIndex[T comparable](collection []T, el T) int {
	for i := range collection {
		if reflect.DeepEqual(collection[i], el) {
			return i
		}
	}
	return -1
}

func Contains[T comparable](b []T, e T) bool {
	for _, v := range b {
		if reflect.DeepEqual(e, v) {
			return true
		}
	}
	return false
}

// DiffOneWay returns the elements in `a` that aren't in `b`.
func DiffOneWay(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func DiffString(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

func AnyToString(slice []any) []string {
	s := make([]string, len(slice))
	for i, v := range slice {
		s[i] = v.(string)
	}
	return s
}
