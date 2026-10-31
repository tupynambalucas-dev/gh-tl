package ui

import (
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/term"
)

var (
	isTerminal = term.IsTerminal(os.Stdout)
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

func colorize(code, text string) string {
	if !isTerminal {
		return text
	}
	return code + text + colorReset
}

func Cyan(text string) string   { return colorize(colorCyan, text) }
func Green(text string) string  { return colorize(colorGreen, text) }
func Red(text string) string    { return colorize(colorRed, text) }
func Yellow(text string) string { return colorize(colorYellow, text) }
func Gray(text string) string   { return colorize(colorGray, text) }
func Bold(text string) string   { return colorize(colorBold, text) }

func PrintHeader(title string) {
	sep := "=========================================================="
	fmt.Println(Cyan(sep))
	fmt.Println(Cyan(title))
	fmt.Println(Cyan(sep))
	fmt.Println()
}

func PrintSubheader(title string) {
	fmt.Println(Bold(title))
}

func PrintPass(msg string) {
	fmt.Println(Green(" PASS ") + msg)
}

func PrintFail(msg string) {
	fmt.Println(Red(" FAIL ") + msg)
}

func PrintWarn(msg string) {
	fmt.Println(Yellow(" WARN ") + msg)
}

func PrintInfo(msg string) {
	fmt.Println(Cyan(" INFO ") + msg)
}

func PrintRemediation(msg string) {
	fmt.Println(Yellow("  Remediation: ") + msg)
}
