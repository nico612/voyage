package term

import (
	"fmt"
	"github.com/moby/term"
	"io"
)

// TerminalSize 获取用户终端的尺寸信息
func TerminalSize(w io.Writer) (int, int, error) {
	outFd, isTerminal := term.GetFdInfo(w)
	if !isTerminal {
		return 0, 0, fmt.Errorf("given writer is no terminal")
	}
	winsize, err := term.GetWinsize(outFd)
	if err != nil {
		return 0, 0, err
	}
	return int(winsize.Width), int(winsize.Height), nil
}
