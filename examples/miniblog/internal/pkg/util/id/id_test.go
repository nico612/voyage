package id

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// 单元测试
func TestGenShortID(t *testing.T) {
	shortID := GenShortID()
	assert.NotEqual(t, "", shortID)
	assert.Equal(t, 6, len(shortID))
}

// 基准测试（性能测试用例） 测试命令：`go test -benchmem -bench .`
// BenchmarkGenShortID-12           2702156               422.7 ns/op            84 B/op          4 allocs/op
// 执行结果分析：
// - 测试函数名称：BenchmarkGenShortID-12
// - 执行了2702156次，
// - 每次的执行平均时间 422.7 纳秒，
// - 84 B/op：每次迭代的内存分配量 84 字节，
// - 4 allocs/op: 每次迭代的内存分配次数为。
func BenchmarkGenShortID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenShortID()
	}
}

func BenchmarkGenShortIDTimeConsuming(b *testing.B) {
	b.StopTimer() //调用该函数停止压力测试的时间计数

	shortId := GenShortID()
	if shortId == "" {
		b.Error("Failed to generate short id")
	}

	// 在 b.StopTimer() 和 b.StartTimer() 之间可以做一些准备工作，这样这些时间不影响我们测试函数本身的性能。

	b.StartTimer() //重新开始时间

	for i := 0; i < b.N; i++ {
		GenShortID()
	}
}

// 查看性能并生成函数调用图 `go test -bench= . -cpuprofile=cpu.profile`
// 上述命令会在当前目录下生成 cpu.profile 和 id.test 文件。
// 之后，我们可以执行 go tool pprof id.test cpu.profile 查看性能（进入交互界面后执行 top 指令）：
