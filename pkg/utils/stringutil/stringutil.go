package stringutil

import (
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
)

// Diff 从一个基础字符串切片中排除另一个排除字符串切片中包含的元素，返回剩余的元素构成的切片。
func Diff(base, exclude []string) (result []string) {
	excludeMap := make(map[string]bool)
	for _, s := range exclude {
		excludeMap[s] = true
	}
	for _, s := range base {
		if !excludeMap[s] {
			result = append(result, s)
		}
	}
	return result
}

// Unique 一个字符串切片中去除重复元素，返回一个包含唯一值的切片。
func Unique(ss []string) (result []string) {
	smap := make(map[string]bool)
	for _, s := range ss {
		smap[s] = true
	}
	for s := range smap {
		result = append(result, s)
	}
	return result
}

// CamelCaseToUnderscore 驼峰转下划线
// Ex.: MyFunc => my_func
func CamelCaseToUnderscore(str string) string {
	return govalidator.CamelCaseToUnderscore(str)
}

// UnderscoreToCamelCase 下划线转驼峰
// Ex.: my_func => MyFunc
func UnderscoreToCamelCase(str string) string {
	return govalidator.UnderscoreToCamelCase(str)
}

// FindString 在给定的字符串切片 array 中查找特定的字符串 str，并返回其在切片中的索引。
func FindString(array []string, str string) int {
	for index, s := range array {
		if str == s {
			return index
		}
	}
	return -1
}

// StringIn 判断一个字符串 str 是否存在于给定的字符串切片 array 中。
func StringIn(str string, array []string) bool {
	return FindString(array, str) > -1
}

// Reverse 字符串 s 反转，并返回反转后的字符串。
func Reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}
