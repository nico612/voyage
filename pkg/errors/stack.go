package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// 该包主要用户定位发生错误的文件路径、方法名、行号，以便排查问题

// Frame 代表堆栈帧内的程序计数器。出于历史原因，如果将Frame解释为uintptr，则其值表示程序计数器加1。
type Frame uintptr

// pc 返回此帧的程序计数器；多个帧可能具有相同的PC值。
func (f Frame) pc() uintptr {
	return uintptr(f) - 1
}

// file 返回包含此帧的PC所在函数的完整文件路径。
func (f Frame) file() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

// line 返回此帧的PC所在函数源代码的行号。
func (f Frame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

// name 返回函数名
func (f Frame) name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

// Format 按照fmt.Formatter接口格式化帧。
//
//	%s    源文件
//	%d    源代码行号
//	%n    函数名称
//	%v    等同于 %s:%d
//
// Format接受一些标志，用于修改某些动词的打印方式，如下：
//
//	%+s   函数名称和源文件路径，相对于编译时的GOPATH，由\n\t分隔 (<funcname>\n\t<path>)
//	%+v   等同于 %+s:%d
func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			io.WriteString(s, f.name())
			io.WriteString(s, "\n\t")
			io.WriteString(s, f.file())
		default:
			io.WriteString(s, path.Base(f.file()))
		}
	case 'd':
		io.WriteString(s, strconv.Itoa(f.line()))
	case 'n':
		io.WriteString(s, funcname(f.name()))
	case 'v':
		f.Format(s, 's')
		io.WriteString(s, ":")
		f.Format(s, 'd')
	}
}

// MarshalText 将堆栈跟踪的 Frame 格式化为文本字符串。输出与fmt.Sprintf("%+v", f)相同，但不包含换行符或制表符。
func (f Frame) MarshalText() ([]byte, error) {
	name := f.name()
	if name == "unknown" {
		return []byte(name), nil
	}
	return []byte(fmt.Sprintf("%s %s:%d", name, f.file(), f.line())), nil
}

// StackTrace 是从最内部（最新）到最外部（最旧）的帧堆栈。
type StackTrace []Frame

// Format 根据fmt.Formatter接口格式化帧堆栈。
//
//	%s	列出堆栈中每个帧的源文件
//	%v	列出堆栈中每个帧的源文件和行号
//
// Format接受一些标志，用于修改某些动词的打印方式，如下：
//
//	%+v   对堆栈中每个帧打印文件名、函数名和行号。
func (st StackTrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			for _, f := range st {
				io.WriteString(s, "\n")
				f.Format(s, verb)
			}
		case s.Flag('#'):
			fmt.Fprintf(s, "%#v", []Frame(st))
		default:
			st.formatSlice(s, verb)
		}
	case 's':
		st.formatSlice(s, verb)
	}
}

// formatSlice 将此 StackTrace 格式化为给定缓冲区的 Frame 切片，在使用'%s'或'%v'时才有效。
func (st StackTrace) formatSlice(s fmt.State, verb rune) {
	io.WriteString(s, "[")
	for i, f := range st {
		if i > 0 {
			io.WriteString(s, " ")
		}
		f.Format(s, verb)
	}
	io.WriteString(s, "]")
}

// stack 代表一组程序计数器。
type stack []uintptr

func (s *stack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('+'):
			for _, pc := range *s {
				f := Frame(pc)
				fmt.Fprintf(st, "\n%+v", f)
			}
		}
	}
}

func (s *stack) StackTrace() StackTrace {
	f := make([]Frame, len(*s))
	for i := 0; i < len(f); i++ {
		f[i] = Frame((*s)[i])
	}
	return f
}

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

// funcname 移除由func.Name()报告的函数名称的路径前缀组件。
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
