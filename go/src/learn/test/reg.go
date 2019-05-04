package main

import (
	"fmt"
	"regexp"
)

const text = `
	mail is sakuraus@gmail.com
mail2 is 1757510543@qq.com
mail3 is 18215675751@163.com.cn aaa
`

/**

. 	任意字符
+    一个/多个
* 	0个/多个
() 	提取的部分
[]	匹配区间
\	转义字符

 */
func main() {
	rule := `([a-zA-Z0-9]+)@([a-z-A-Z0-9]+)\.([a-z-A-Z0-9.]+)`
	compile := regexp.MustCompile(rule)

	submatch := compile.FindAllStringSubmatch(text, -1)
	for _, match := range submatch {
		fmt.Println(match)
	}
}
