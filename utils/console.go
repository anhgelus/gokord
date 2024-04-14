package utils

import (
	"fmt"
	"log"
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

var DebugEnabled bool

// SendSuccess sends a success message
func SendSuccess(message string) {
	log.Default().Println(AnsiGreen, message, AnsiReset)
}

// SendWarn sends a warning message
func SendWarn(message string) {
	log.Default().Println(AnsiYellow, message, AnsiReset)
}

// SendDebug sends a debug message
func SendDebug(message ...any) {
	if DebugEnabled {
		log.Default().Println(AnsiCyan, message, AnsiReset)
	}
}

// SendAlert sends an alert
func SendAlert(pos string, message string) {
	log.Default().Println(fmt.Sprintf("[%s] %s%s%s", pos, AnsiRed, message, AnsiReset))
}

// SendError sends an error (like a panic(any...))
func SendError(err error) {
	log.Default().Panic(err)
}
