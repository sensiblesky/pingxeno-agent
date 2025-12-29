//go:build !windows
// +build !windows

package main

import "fmt"

// runGUI is not available on non-Windows platforms
func runGUI() error {
	return fmt.Errorf("GUI mode is only available on Windows")
}

