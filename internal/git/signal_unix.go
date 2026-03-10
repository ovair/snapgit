//go:build !windows

package git

import (
	"os"
	"syscall"
)

var forwardedSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
