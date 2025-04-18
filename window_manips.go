package main

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

type RECT struct {
	Left, Top, Right, Bottom int32
}

func FindWindow(title string) syscall.Handle {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	ret, _, _ := procFindWindow.Call(
		0,
		uintptr(unsafe.Pointer(titlePtr)),
	)
	return syscall.Handle(ret)
}

func GetWindowRect(hwnd syscall.Handle) (RECT, error) {
	var rect RECT
	ret, _, _ := procGetWindowRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&rect)),
	)
	if ret == 0 {
		return rect, fmt.Errorf("GetWindowRect failed")
	}
	return rect, nil
}

func SetForegroundWindow(hwnd syscall.Handle) bool {
	ret, _, _ := procSetForegroundWindow.Call(uintptr(hwnd))
	return ret != 0
}

func WaitForWindow(title string, timeout time.Duration) (syscall.Handle, error) {
	start := time.Now()
	for time.Since(start) < timeout {
		hwnd := FindWindow(title)
		if hwnd != 0 {
			return hwnd, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return 0, fmt.Errorf("window not found after %v", timeout)
}
