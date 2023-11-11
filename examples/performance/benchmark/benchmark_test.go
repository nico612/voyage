package benchmark

import (
	"math"
	"testing"
)

func BenchmarkExp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = math.Exp(3.5)
	}
}

// 执行命令。采集性能数据 `go test -bench=BenchmarkExp -benchmem -cpuprofile cpu.profile -memprofile mem.profile`
