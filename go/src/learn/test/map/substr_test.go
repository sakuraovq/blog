package _map

import "testing"

//  go test .
// golang中表格测试
func TestCuntSubStr(t *testing.T) {
	tests := []struct {
		s   string
		ans int
	}{
		{"abc", 3},
		// chinese support
		{"在哦这里哦哦", 4},

		// other
		{"bbb", 1},
	}

	for _, test := range tests {
		result := countSubStr(test.s)
		if result != test.ans {
			t.Errorf("countSubstr s=%s ans=%d expected=%d",
				test.s, result, test.ans)
		}

	}
}

//  go test -bench .  #性能测试

// go test -bench . cpuprofile cpu.out
// go tool pprof cpu.out
// web  #需要安装svg工具
func BenchmarkSubStr(t *testing.B) {
	s := "黑化肥挥发黑"
	ans := 5

	for i := 0; i < t.N; i++ {
		result := countSubStr(s)
		if result != ans {
			t.Errorf("countSubstr s=%s ans=%d expected=%d",
				s, result, ans)
		}
	}

}
