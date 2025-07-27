//go:build !windows

package main

func startMSI(msiPath string) error {
	// Linux/macOS 上无 MSI，留空或做提示
	return nil
}
