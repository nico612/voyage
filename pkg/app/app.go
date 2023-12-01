package app

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	flag2 "github.com/nico612/go-project/pkg/cli/flag"
	"github.com/nico612/go-project/pkg/cli/globalflag"
	"github.com/nico612/go-project/pkg/errors"
	"github.com/nico612/go-project/pkg/log"
	"github.com/nico612/go-project/pkg/term"
	"github.com/nico612/go-project/pkg/version"
	"github.com/nico612/go-project/pkg/version/verflag"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

// 构建 App 命令行

var (
	progressMessage = color.GreenString("==>")
)

type App struct {
	basename    string
	name        string
	description string

	options CliOptions // 标志(flag)集
	runFunc RunFunc

	noVersion bool // 不提供版本标志
	noConfig  bool // 不提供配置标志
	silence   bool // 错误时打印使用信息

	commands []*Command // 子命令
	args     cobra.PositionalArgs
	cmd      *cobra.Command
}

type Option func(*App)

// WithOptions 设置标志集
func WithOptions(opt CliOptions) Option {
	return func(app *App) {
		app.options = opt
	}
}

// RunFunc defines the application's startup callback function.
type RunFunc func(basename string) error

func WithRunFunc(run RunFunc) Option {
	return func(app *App) {
		app.runFunc = run
	}
}

// WithDescription 命令行描述
func WithDescription(desc string) Option {
	return func(app *App) {
		app.description = desc
	}
}

func WithSilence() Option {
	return func(app *App) {
		app.silence = true
	}
}

// WithNoVersion 应用程序将不提供版本（version）flag
func WithNoVersion() Option {
	return func(app *App) {
		app.noVersion = true
	}
}

// WithNoConfig 不提供 config 标志
func WithNoConfig() Option {
	return func(app *App) {
		app.noConfig = true
	}
}

// WithValidArgs 选项参数
func WithValidArgs(args cobra.PositionalArgs) Option {
	return func(app *App) {
		app.args = args
	}
}

// WithDefaultValidArgs 设置默认的验证函数，用于验证非标志参数的有效性。在这个函数中，
// 通过设置 a.args 字段为一个闭包函数来实现。
func WithDefaultValidArgs() Option {
	return func(a *App) {
		a.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {

				// 处理非标志参数逻辑
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}

			}
			return nil
		}
	}
}

func NewApp(name string, basename string, opts ...Option) *App {
	a := &App{
		basename: basename,
		name:     name,
	}

	for _, opt := range opts {
		opt(a)
	}

	a.buildCommand()

	return a
}

func (a *App) buildCommand() {
	cmd := cobra.Command{
		Use:           a.basename,
		Short:         a.name,
		Long:          a.description,
		Args:          a.args,
		SilenceUsage:  true, // 执行错误时 显示帮助信息
		SilenceErrors: true, // 执行错误时退出
	}

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true

	//将 Go 标准库的命令行标志集成到 Cobra 应用中。
	cmd.Flags().AddGoFlagSet(flag.CommandLine)

	// 添加子命令
	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		cmd.SetHelpCommand(helpCommand(FormatBaseName(a.basename)))
	}

	// 命令运行函数
	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}

	// 添加标志集
	var namedFlagSets flag2.NamedFlagSets
	if a.options != nil {
		namedFlagSets = a.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f)
		}
	}

	if !a.noVersion { // 添加版本标志
		verflag.AddFlags(namedFlagSets.FlagSet("global"))
	}

	if !a.noConfig { // 添加config标志
		addConfigFlag(a.basename, namedFlagSets.FlagSet("global"))
	}

	globalflag.AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name())

	// 将 global 中的所有标志添加到 cobra 中
	cmd.Flags().AddFlagSet(namedFlagSets.FlagSet("global"))

	addCmdTemplate(&cmd, namedFlagSets)

	a.cmd = &cmd
}

func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		log.Infof("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

func (a *App) Command() *cobra.Command {
	return a.cmd
}

func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	printWorkingDir()
	flag2.PrintFlags(cmd.Flags())

	if !a.noVersion { // 打印版本并退出
		verflag.PrintAndExitIfRequested()
	}

	if !a.noConfig { // 解析配置文件

		// 通过将 Cobra 命令的标志与 Viper 绑定，可以让 Viper 管理这些标志的值，实现了标志值的读取和更新，
		// 并且允许从多个来源（例如命令行、环境变量、配置文件等）获取和设置标志的值。
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}

		// 解析配置文件到指定结构体
		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}

	if !a.silence {
		log.Infof("%v Starting %s ...", progressMessage, a.name)
		if !a.noVersion {
			log.Infof("%v Version: `%s`", progressMessage, version.Get().ToJSON())
		}

		if !a.noConfig {
			log.Infof("%v Config file used: `%s`", progressMessage, viper.ConfigFileUsed())
		}
	}

	// 处理各项配置的规则
	if a.options != nil {
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}

	// run application
	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}

	return nil
}

func (a *App) applyOptionRules() error {

	// 检查各项配置是否准备完成
	if completeableOptions, ok := a.options.(CompleteableOptions); ok {
		if err := completeableOptions.Complete(); err != nil {
			return err
		}
	}

	// 各项配置验证
	if errs := a.options.Validate(); len(errs) != 0 {
		return errors.NewAggregate(errs)
	}

	// 配置打印
	if printableOptions, ok := a.options.(PrintableOptions); ok && !a.silence {
		log.Infof("%v Config: `%s`", progressMessage, printableOptions.String())
	}

	return nil
}

func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Infof("%v WorkingDir: %s", progressMessage, wd)
}

func printFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		log.Infof("FLAG: --%s=%q", flag.Name, flag.Value)
	})
}

func addCmdTemplate(cmd *cobra.Command, namedFlagSets flag2.NamedFlagSets) {
	usageFmt := "Usage:\n  %s\n"
	// 获取终端尺寸信息， 以便格式化输出
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())

	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		flag2.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)

		return nil
	})

	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		flag2.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})
}
