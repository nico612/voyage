package sliceutil

// RemoveString 从字符串切片中移除满足特定条件的字符串元素，并返回处理后的字符串切片。
func RemoveString(slice []string, remove func(item string) bool) []string {
	for i := 0; i < len(slice); i++ {
		if remove(slice[i]) {
			slice = append(slice[:i], slice[i+1:]...)
			i--
		}
	}
	return slice
}

// FindString 在字符串切片中查找目标字符串 target 是否存在
func FindString(slice []string, target string) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}

// FindInt 在整数切片中查找目标整数 target 是否存在。
func FindInt(slice []int, target int) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}

// FindUint 在uint切片中查找目标 target 是否存在。
func FindUint(slice []uint, target uint) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}
