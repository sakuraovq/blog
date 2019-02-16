package _map

import "fmt"

func main() {

	fmt.Println(countSubStr("哈哈是哈哈我"))
}

// 寻找不重复最大子串的长度
func countSubStr(s string) int {
	lastCurrent := make(map[rune]int)
	start := 0
	maxLenth := 0
	for i, char := range []rune(s) {
		// 如果 lastI出现在 map字典的后面则 start 向后移动,小于start或者不存在 不重复子串长度不受影响
		if lastI, ok := lastCurrent[char]; ok && lastI >= start {
			start = lastI + 1
		}
		if i-start+1 > maxLenth {
			maxLenth = i - start + 1
		}
		lastCurrent[char] = i
	}
	return maxLenth
}
