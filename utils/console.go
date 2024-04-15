package utils

import (
	"log/slog"
)

const (
	AnsiReset       = "\033[0m"
	AnsiRed         = "\033[31m"
	AnsiGreen       = "\033[32m"
	AnsiYellow      = "\033[33m"
	AnsiBlue        = "\033[34m"
	AnsiMagenta     = "\033[35m"
	AnsiCyan        = "\033[36m"
	AnsiWhite       = "\033[37m"
	AnsiBlueBold    = "\033[34;1m"
	AnsiMagentaBold = "\033[35;1m"
	AnsiRedBold     = "\033[31;1m"
	AnsiYellowBold  = "\033[33;1m"
)

var DebugEnabled = true

// SendSuccess sends a success message
func SendSuccess(message string, args ...any) {
	slog.Info(AnsiGreen+message+AnsiReset, args)
}

// SendWarn sends a warning message
func SendWarn(message string, args ...any) {
	slog.Warn(AnsiYellow+message+AnsiReset, args)
}

// SendDebug sends a debug message
func SendDebug(message string, args ...any) {
	if DebugEnabled {
		slog.Debug(AnsiCyan+message+AnsiReset, args)
	}
}

// SendAlert sends an alert
func SendAlert(pos string, message string) {
	slog.Error(AnsiRed+message+AnsiReset, "position", pos)
}
