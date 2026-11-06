//go:build fyne && windows

package cmd

import (
	"syscall"
	"unsafe"
)

// The window executable is linked without a console (-H windowsgui), so an
// error printed to stderr during startup would vanish. A message box is the one
// surface Windows guarantees, even before Fyne has an app or a window.
var (
	user32DLL  = syscall.NewLazyDLL("user32.dll")
	messageBox = user32DLL.NewProc("MessageBoxW")
)

const messageBoxIconError = 0x00000010

func showFyneStartupError(message string) {
	text, err := syscall.UTF16PtrFromString(message)
	if err != nil {
		return
	}

	title, err := syscall.UTF16PtrFromString("crower")
	if err != nil {
		return
	}

	messageBox.Call(0, uintptr(unsafe.Pointer(text)), uintptr(unsafe.Pointer(title)), messageBoxIconError)
}
