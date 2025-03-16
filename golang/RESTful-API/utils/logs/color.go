package logs

import (
	"fmt"
	"strings"
)

const (
	colorRed = uint8(iota + 91)
	colorGreen
	colorYellow
	colorBlue
	colorMagenta //洋红
)

// 根据不同日志级别添加颜色
func ColorLevel(logLevel string) string {
	var (
		colorString = strings.ToLower(logLevel)
	)
	switch colorString {
	case "info":
		colorString = InfoColor(colorString)
	case "warn", "warning":
		colorString = WarnColor(colorString)
	case "debug":
		colorString = SuccessColor(colorString)
	case "trace":
		colorString = TraceColor(colorString)
	case "fatal", "panic", "error":
		colorString = ErrorColor(colorString)
	}
	return colorString
}

func TraceColor(format string) string {
	return yellow(format)
}
func InfoColor(format string) string {
	return blue(format)
}
func WarnColor(format string) string {
	return magenta(format)
}
func SuccessColor(format string) string {
	return green(format)
}
func ErrorColor(format string) string {
	return red(format)
}
func red(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorRed, s)
}
func green(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorGreen, s)
}
func yellow(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorYellow, s)
}
func blue(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorBlue, s)
}
func magenta(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorMagenta, s)
}
