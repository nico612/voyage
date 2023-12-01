package flag

import (
	goflag "flag"
	"fmt"
	"github.com/spf13/pflag"
	"strings"
)

// WordSepNormalizeFunc changes all flags that contain "_" separators. 将所有的包含 "_" 分割符 的标志 改为使用 "-" 分割符
func WordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
	}
	return pflag.NormalizedName(name)
}

// WarnWordSepNormalizeFunc changes and warns for flags that contain "_" separators.
func WarnWordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		nname := strings.ReplaceAll(name, "_", "-")

		fmt.Printf("%s is DEPRECATED and will be removed in a future version. Use %s instead.", name, nname)

		return pflag.NormalizedName(nname)
	}
	return pflag.NormalizedName(name)
}

// InitFlags normalizes, parses, then logs the command line flags.
func InitFlags(flags *pflag.FlagSet) {
	// 设置标志命名的规则，单词间使用 "-" 分割符连接
	flags.SetNormalizeFunc(WordSepNormalizeFunc)

	// 将 Go 标准库的命令行标志集成到 Cobra 应用中。
	flags.AddGoFlagSet(goflag.CommandLine)
}

func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		fmt.Printf("FLAG: --%s=%q", flag.Name, flag.Value)
	})
}
