package test

import (
	"fmt"
	"os"
	"testing"
)

// go test ：运行改包所有的单元测试 -v 显示详细结果 -cover 查看测试覆盖率
// go test -run TestAdd -v 指定运行 TestAdd 测试用例
func TestAdd(t *testing.T) {
	if ans := Add(1, 2); ans != 3 {
		t.Errorf("1 + 2 expected be 3, but %d got", ans)
	}

	if ans := Add(-10, -20); ans != -30 {
		t.Errorf("-10 + -20 expected be -30, but %d got", ans)
	}
}

// 子测试
func TestMul(t *testing.T) {

	// 子测试，在某个测试用例中，根据测试场景使用 t.Run创建不同的子测试用例：
	t.Run("pos", func(t *testing.T) {
		if Mul(2, 3) != 6 {
			t.Fatal("fail")
		}
	})

	t.Run("neg", func(t *testing.T) {
		if Mul(2, -3) != -6 {
			t.Fatal("fail")
		}
	})
}

// 多个子测试的场景，更推荐如下的写法(table-driven tests)：
// 所有用例的数据组织在切片 cases 中，看起来就像一张表，借助循环创建子测试。这样写的好处有：
//
// 新增用例非常简单，只需给 cases 新增一条测试数据即可。
// 测试代码可读性好，直观地能够看到每个子测试的参数和期待的返回值。
// 用例失败时，报错信息的格式比较统一，测试报告易于阅读。
// 如果数据量较大，或是一些二进制数据，推荐使用相对路径从文件中读取。
func TestMul2(t *testing.T) {
	cases := []struct {
		Name           string
		A, B, Expected int
	}{
		{"pos", 2, 3, 6},
		{"neg", 2, -3, -6},
		{"zero", 2, 0, 0},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if ans := Mul(c.A, c.B); ans != c.Expected {
				t.Fatalf("%d * %d expected %d, but %d got",
					c.A, c.B, c.Expected, ans)
			}
		})
	}
}

// 帮助函数
//对一些重复的逻辑，抽取出来作为公共的帮助函数(helpers)，可以增加测试代码的可读性和可维护性。 借助帮助函数，可以让测试用例的主逻辑看起来更清晰。

type calcCase struct {
	A, B, Expected int
}

func createMulTestCase(t *testing.T, c *calcCase) {
	t.Helper() //用于标注该函数是帮助函数，报错时将输出帮助函数 调用者 的信息(文件和行号)，而不是帮助函数的内部信息。
	if ans := Mul(c.A, c.B); ans != c.Expected {
		t.Fatalf("%d * %d expected %d, but %d got",
			c.A, c.B, c.Expected, ans)
	}
}
func TestMul3(t *testing.T) {
	createMulTestCase(t, &calcCase{2, 3, 6})
	createMulTestCase(t, &calcCase{2, -3, -6})
	createMulTestCase(t, &calcCase{2, 0, 1}) // wrong case
}

func setup() {
	fmt.Println("Before all tests")
}

func teardown() {
	fmt.Println("After all tests")
}

func Test1(t *testing.T) {
	fmt.Println("I'm test1")
}

func Test2(t *testing.T) {
	fmt.Println("I'm test2")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
