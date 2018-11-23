package main

import "fmt"

func main() {

	// 切片是半开 闭区间截取 cap 是数组的容量 截取超过cap会报错,可以使用append 对切片扩容
	array := [...]int{0, 1, 2, 3, 4, 5, 6, 7}
	fmt.Println("pop")
	end := array[len(array)-1]

	fmt.Println(end)
	s1 := array[:len(array)-1]
	printerSlice(s1)

	fmt.Println("shift")
	head := array[0]
	s2 := array[1:]
	fmt.Println(head)
	printerSlice(s2)

	fmt.Println("delete middle")
	// 从前面取2个舍弃第三个元素, 从第四个元素开始取
	s3 := append(array[:3], array[4:]...)
	printerSlice(s3)
	// 预生成silce 长度10 底层数组长度15 扩容cap golang是*2
	preSlice := make([]int, 10, 15)
	preSlice[1] = 2
	printerSlice(preSlice)

	// 空的 slice不能直接赋值需要make因为底层没有数据 , slice只是对数据的 一个view
	var nilSlice []int
	//nilSlice[0] = 1
	printerSlice(nilSlice)
}

func printerSlice(s []int) {

	fmt.Println(s)
	fmt.Printf(" len=%d cap=%d \n", len(s), cap(s))
}
