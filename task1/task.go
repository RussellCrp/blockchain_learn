package main

import (
	"container/list"
	"fmt"
	"sort"
)

// 只出现一次的数字：
// 给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
func singleNumber(nums []int) int {
	mapping := make(map[int]int)
	for _, val := range nums {
		mapping[val] = mapping[val] + 1
	}

	for key, val := range mapping {
		if val == 1 {
			return key
		}
	}
	return -1
}

// 有效的括号
// 给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。
func isValid(s string) bool {
	m := make(map[rune]rune, 3)
	m['('] = ')'
	m['{'] = '}'
	m['['] = ']'
	l := list.New()
	for _, val := range s {
		_, ok := m[val]
		if ok {
			l.PushFront(val)
		} else {
			openChar := l.Front()
			if openChar == nil {
				return false
			}
			if m[openChar.Value.(rune)] == val {
				l.Remove(openChar)
			} else {
				return false
			}
		}
	}
	return l.Len() == 0
}

// 编写一个函数来查找字符串数组中的最长公共前缀。
// 如果不存在公共前缀，返回空字符串 ""。
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	maxPrefix := strs[0]
	for i := 1; i < len(strs); i++ {
		str := strs[i]
		var length int = len(maxPrefix)
		if len(str) < len(maxPrefix) {
			length = len(str)
		}
		var ci int
		for {
			if ci >= length {
				break
			}
			if str[ci] == maxPrefix[ci] {
				ci++
			} else {
				break
			}
		}
		maxPrefix = maxPrefix[:ci]
	}
	return maxPrefix
}

// 加一
// 给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一
func plusOne(digits []int) []int {
	index, length := 1, len(digits)
	for {
		if digits[length-index] == 9 {
			digits[length-index] = 0
			index++
			if index > length {
				digits = append([]int{1}, digits...)
				break
			}
		} else {
			digits[length-index] = digits[length-index] + 1
			break
		}
	}
	return digits
}

// 删除有序数组中的重复项
func removeDuplicates(nums []int) int {
	slow := 0
	for fast := 1; fast < len(nums); fast++ {
		if nums[slow] != nums[fast] {
			slow++
			nums[slow] = nums[fast]
		}
	}
	return slow + 1
}

// 合并区间
func merge(intervals [][]int) [][]int {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	result := make([][]int, 0, len(intervals))
	result = append(result, intervals[0])
	for i, j := 0, 1; j < len(intervals); j++ {
		resultI := result[i]
		interval := intervals[j]
		if resultI[1] >= interval[0] {
			maxNum := interval[1]
			if resultI[1] > interval[1] {
				maxNum = resultI[1]
			}
			result[i] = []int{resultI[0], maxNum}
		} else {
			i++
			result = append(result, intervals[j])
		}
	}
	return result
}

func main() {
	fmt.Println(singleNumber([]int{2, 2, 1}))
	fmt.Println(isValid("()"))
	fmt.Println(longestCommonPrefix([]string{"flower", "flow", "flight"}))
	fmt.Println(plusOne([]int{9, 9, 9, 9}))
	ints := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	fmt.Println(ints[:removeDuplicates(ints)])
	fmt.Println(merge([][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}))
}
