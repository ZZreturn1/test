package common

import "sort"

func IsSubString(target string, str_array []string) bool {
               // 对字符串切片进行排序
	sort.Strings(str_array) 
                // 使用二分查找在排序后的切片中查找目标字符串的位置
	index := sort.SearchStrings(str_array, target) 
                // 判断目标字符串是否在切片中找到
	return index < len(str_array) && str_array[index] == target 
}