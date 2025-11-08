package commonutil

import (
	"sort"
	"strings"
)

func InSlice(str string, list []string) bool {
	for _, s := range list {
		if str == s {
			return true
		}
	}

	return false
}

func InIntSlice(str int, list []int) bool {
	for _, s := range list {
		if str == s {
			return true
		}
	}

	return false
}

func UniqueSlice(list []string) []string {
	m := make(map[string]bool)
	uniqueList := []string{}
	for _, e := range list {
		if !m[e] {
			uniqueList = append(uniqueList, e)
			m[e] = true
		}
	}
	return uniqueList
}

func SliceIntersect(listA, listB []string) []string {
	res := []string{}
	if len(listA) == 0 || len(listB) == 0 {
		return res
	}
	m := make(map[string]bool)
	for _, e := range listA {
		m[e] = true
	}
	for _, e := range listB {
		if m[e] {
			res = append(res, e)
		}
	}
	return res
}

func SliceOverlap(listA, listB []string) bool {
	if len(listA) == 0 || len(listB) == 0 {
		return false
	}
	m := make(map[string]bool)
	for _, e := range listA {
		m[e] = true
	}
	for _, e := range listB {
		if m[e] {
			return true
		}
	}
	return false
}

func SliceElementsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	if len(a) == 0 {
		return true
	}

	m := make(map[string]int, len(a))
	for i := range a {
		m[a[i]]++
	}

	for i := range b {
		if v, ok := m[b[i]]; !ok || v == 0 {
			return false
		} else {
			m[b[i]]--
		}
	}
	return true
}

func SliceEqualRegardlessOfOrder(listA, listB []string) bool {
	if len(listA) != len(listB) {
		return false
	}
	a, b := make([]string, len(listA)), make([]string, len(listB))
	copy(a, listA)
	copy(b, listB)
	sort.Strings(a)
	sort.Strings(b)
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func SliceSubstract(listA, listB []string) []string {
	res := []string{}
	m := make(map[string]bool)
	for _, e := range listB {
		m[e] = true
	}
	for _, e := range listA {
		if !m[e] {
			res = append(res, e)
		}
	}
	return res
}

func ConvertMapToStrings(mapLists map[string][]string) []string {
	list := []string{}
	for _, val := range mapLists {
		list = append(list, val...)
	}
	return list
}

// ConvertSliceToSlice2D is used to convert slice to two-dimensional slice.
func ConvertSliceToSlice2D(list []string, num int) [][]string {
	count := len(list) / num
	remainder := len(list) % num

	// calculate silce capacity
	capacity := count
	if remainder > 0 {
		capacity++
	}

	var result [][]string = make([][]string, capacity)
	var i int = 0
	for ; i < count; i++ {
		result[i] = list[i*num : i*num+num]
	}

	if remainder > 0 {
		result[i] = list[i*num:]
	}

	return result
}

func SplitByMultiSepAndTrimSpace(ori string, ss ...string) []string {
	if len(ss) == 0 {
		return []string{strings.TrimSpace(ori)}
	}
	f := ss[0]
	for _, s := range ss[1:] {
		ori = strings.Replace(ori, s, f, -1)
	}
	var res []string
	for _, s := range strings.Split(ori, f) {
		if m := strings.TrimSpace(s); m != "" {
			res = append(res, m)
		}
	}
	return res
}
