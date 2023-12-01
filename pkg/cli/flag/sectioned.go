package flag

import (
	"bytes"
	"fmt"
	"github.com/spf13/pflag"
	"io"
	"strings"
)

// NamedFlagSets 按照调用 FlagSet 的顺序存储命名标志集。
type NamedFlagSets struct {
	// Order 标志集 名称 的有序列表。
	Order []string

	// FlagSets 通过 name 储存标志集
	FlagSets map[string]*pflag.FlagSet
}

// FlagSet 获取标志集
func (nfs *NamedFlagSets) FlagSet(name string) *pflag.FlagSet {
	if nfs.FlagSets == nil {
		nfs.FlagSets = map[string]*pflag.FlagSet{}
	}

	if _, ok := nfs.FlagSets[name]; !ok {
		//新建一个命令标志，遇到错误时直接退出，不进行错误处理。
		nfs.FlagSets[name] = pflag.NewFlagSet(name, pflag.ExitOnError)
		nfs.Order = append(nfs.Order, name)
	}

	return nfs.FlagSets[name]
}

// PrintSections 根据指定的列数打印出一组命令行标志的使用说明，并根据特定的条件对输出的格式进行处理和调整。 如果 cols 为零，则不换行。
func PrintSections(w io.Writer, fss NamedFlagSets, cols int) {
	for _, name := range fss.Order {
		fs := fss.FlagSets[name]
		if !fs.HasFlags() {
			continue
		}

		wideFS := pflag.NewFlagSet("", pflag.ExitOnError)
		wideFS.AddFlagSet(fs)

		var zzz string
		if cols > 24 {
			zzz = strings.Repeat("z", cols-24)
			wideFS.Int(zzz, 0, strings.Repeat("z", cols-24))
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "\n%s flags:\n\n%s", strings.ToUpper(name[:1])+name[1:], wideFS.FlagUsagesWrapped(cols))

		if cols > 24 {
			i := strings.Index(buf.String(), zzz)
			lines := strings.Split(buf.String()[:i], "\n")
			fmt.Fprint(w, strings.Join(lines[:len(lines)-1], "\n"))
			fmt.Fprintln(w)
		} else {
			fmt.Fprint(w, buf.String())
		}
	}
}
