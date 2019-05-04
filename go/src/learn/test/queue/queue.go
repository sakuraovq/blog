package queue

type Queue []interface{}

func (q *Queue) Push(value interface{}) {
	*q = append(*q, value.(int))
}

// 使用 可以使用 switch v.(type) 也可以使用 value.(type) 这个两种方式
// 用上了interface{} 运行时才会报错
func (q *Queue) Shift() interface{} {
	if q.Empty() {
		return 0
	}
	head := (*q)[0]
	*q = (*q)[1:]
	return head
}

func (q *Queue) Pop() interface{} {
	if q.Empty() {
		return 0
	}
	end := (*q)[len(*q)-1]

	*q = (*q)[:len(*q)-1]
	return end

}

func (q *Queue) Empty() bool {
	return len(*q) == 0
}
