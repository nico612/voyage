package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"runtime"
	"strings"
)

// Command 自定义命令行
type Command struct {
	usage    string
	short    string
	desc     string
	options  CliOptions     // 提供标志集（flags）的接口
	commands []*Command     // 子命令
	runFunc  RunCommandFunc // 命令行启动回调函数
}

// CommandOption 使用选项模式构建 Command
type CommandOption func(*Command)

// WithCommandOptions 设置命令行配置
// opt：接口类型，该接口提供了 命令标志集和验证的方法
func WithCommandOptions(opt CliOptions) CommandOption {
	return func(command *Command) {
		command.options = opt
	}
}

// RunCommandFunc 定义应用命令行启动回调函数
type RunCommandFunc func(args []string) error

func WithCommandRunFunc(run RunCommandFunc) CommandOption {
	return func(command *Command) {
		command.runFunc = run
	}
}

// NewCommand 新建命令 usage: 使用方法， short: 短描述 desc: 长描述
func NewCommand(usage string, short string, desc string, opts ...CommandOption) *Command {

	c := &Command{
		usage: usage,
		short: short,
		desc:  desc,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// 生成 cobra 的命令
func (c *Command) cobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.usage,
		Short: c.short,
		Long:  c.desc,
	}

	cmd.SetOut(os.Stdout)
	cmd.Flags().SortFlags = false

	// 1. 添加子命令
	if len(c.commands) > 0 {
		for _, command := range c.commands {
			cmd.AddCommand(command.cobraCommand())
		}
	}

	// 2. 命令运行函数
	if c.runFunc != nil {
		cmd.Run = c.runCommand
	}

	// 3. 处理标志(flag)集
	if c.options != nil {
		for _, f := range c.options.Flags().FlagSets {
			cmd.Flags().AddFlagSet(f) // 添加到 cobra 标志集中
		}
	}

	// 4. 添加帮助 flag
	addHelpCommandFlag(c.usage, cmd.Flags())

	return cmd
}

// 命令运行函数
func (c *Command) runCommand(cmd *cobra.Command, args []string) {
	if c.runFunc != nil {
		if err := c.runFunc(args); err != nil {
			fmt.Printf("%v %v\n", color.RedString("Error:"), err)
			os.Exit(1)
		}
	}
}

func (c *Command) AddCommand(cmd *Command) {
	c.commands = append(c.commands, cmd)
}

func (c *Command) AddCommands(cmds ...*Command) {
	c.commands = append(c.commands, cmds...)
}

// FormatBaseName is formatted as an executable file name under different
// operating systems according to the given name.
func FormatBaseName(basename string) string {
	// Make case-insensitive and strip executable suffix if present
	if runtime.GOOS == "windows" {
		basename = strings.ToLower(basename)
		basename = strings.TrimSuffix(basename, ".exe")
	}

	return basename
}
