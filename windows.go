//go:build windows

package main

import "golang.org/x/sys/windows"

func init() {
	// Try to enable VT processing if available.
	// So that colors work on modern win10 if ran through conhost.
	handle, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	if err != nil {
		return
	}

	var mode uint32
	err = windows.GetConsoleMode(handle, &mode)
	if err != nil {
		return
	}

	// ENABLE_VIRTUAL_TERMINAL_PROCESSING
	windows.SetConsoleMode(handle, mode | 0x0004)
}