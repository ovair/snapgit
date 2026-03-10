//go:build windows

package git

import "os"

var forwardedSignals = []os.Signal{os.Interrupt}
