package main

import (
	"fmt"
	"os"
)

func printMaze(maze [][]int) {
	for _, row := range maze {
		for _, val := range row {
			fmt.Printf("%3d", val)
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	maze := readMaze("./maze.in")
	printMaze(maze)
	steps := walk(maze, point{0, 0}, point{len(maze) - 1, len(maze[0]) - 1})
	printMaze(steps)
}

// 数组i 是往下遍历,j 是往右遍历
type point struct {
	i, j int
}

var dirs = [4]point{
	// 搜索顺序 上->左->下->右
	/**
		   i
	       |0  1  2
	  j ---|------
	     0 |[0][0]
		 1 |[0][0]
	 */
	{-1, 0}, {0, -1}, {1, 0}, {0, 1},
}

// 坐标相加
func (p point) add(r point) point {
	return point{p.i + r.i, p.j + r.j}
}


func (p point) at(grid [][]int) (int, bool) {
	// 往上  || 往下 越界  i,j是下标所以统计长度-1
	if p.i < 0 || p.i >= len(grid) {
		return 0, false
	}
	// 往左 || 往右 越界
	if p.j < 0 || p.j >= len(grid[p.i]) {
		return 0, false
	}
	return grid[p.i][p.j], true
}

// start go maze
func walk(maze [][]int, start, end point) [][]int {
	// 记录走过的路~
	steps := make([][]int, len(maze))
	for i := range steps {
		steps[i] = make([]int, len(maze[i]))
	}
	Queue := []point{start}

	for len(Queue) > 0 {
		// 每次消费弹出一个队列头部
		cur := Queue[0]
		Queue = Queue[1:]
		if cur == end {
			break
		}
		for _, dir := range dirs {
			next := cur.add(dir)
			// 迷宫下一个节点是可以走的
			// 并且之前没有走过
			// 和 下一个不等于开始
			val, ok := next.at(maze)
			// 越界了或者撞墙了
			if !ok || val == 1 {
				continue
			}

			val, ok = next.at(steps)
			// 判断是否已经走过了
			if !ok || val != 0 {
				continue
			}

			// 判断下一个节点是否是开始
			if next == start {
				continue
			}
			// 当前步数加1
			curSteps, _ := cur.at(steps)
			steps[next.i][next.j] = curSteps + 1
			// 把可以探索的节点加入队列消费
			Queue = append(Queue, next)
		}
	}
	return steps
}

// 读取迷宫
func readMaze(fileName string) [][]int {

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	var row, col int
	fmt.Fscanf(file, "%d %d", &row, &col)
	maze := make([][]int, row)
	for i := range maze {
		maze[i] = make([]int, col)
		for j := range maze[i] {
			fmt.Fscanf(file, "%d", &maze[i][j])
		}
	}
	return maze
}
