package utils

import "sort"

func SortStrings(slice []string) []string {
	sort.Strings(slice)
	return slice
}
