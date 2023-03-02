//go:build !windows

package cmd

import (
	"fmt"
	"os"
	"runtime"
)

// Unix-specific code for self-updating lvl.
//
// On Unix platforms, we can just move the file on top of the old executable while it's running.
// Easy!

func updateSwapFile(new string, old string) error {
	return os.Rename(new, old)
}

func getAssetFileName() string {
	return fmt.Sprintf("lvl-%s-%s.exe", runtime.GOOS, runtime.GOARCH)
}
