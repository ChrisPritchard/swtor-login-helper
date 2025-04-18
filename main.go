package main

import (
	"fmt"
	"strconv"
	"syscall"
	"time"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")

	// Window functions
	procFindWindow          = user32.NewProc("FindWindowW")
	procGetWindowRect       = user32.NewProc("GetWindowRect")
	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")

	// Input functions
	procSetCursorPos     = user32.NewProc("SetCursorPos")
	procMouseEvent       = user32.NewProc("mouse_event")
	procKeybdEvent       = user32.NewProc("keybd_event")
	procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")
)

const (
	MOUSEEVENTF_LEFTDOWN  = 0x0002
	MOUSEEVENTF_LEFTUP    = 0x0004
	VK_RETURN             = 0x0D
	VK_TAB                = 0x09
	KEYEVENTF_KEYDOWN     = 0x0000
	KEYEVENTF_KEYUP       = 0x0002
	KEYEVENTF_EXTENDEDKEY = 0x0001
	VK_SHIFT              = 0x10
	VK_CONTROL            = 0x11
	VK_MENU               = 0x12 // ALT key
	VK_CAPITAL            = 0x14
)

func main() {

	const windowTitle = "Star Wars™: The Old Republic™"

	// Relative positions within the window
	// These are percentages of window width/height from top-left corner

	// usernameFieldPos := struct{ X, Y float32 }{0.05, 0.4} // note: assuming 'save account' is checked, the username (and only the username) will be prefilled
	passwordFieldPos := struct{ X, Y float32 }{0.05, 0.56}
	otpFieldPos := struct{ X, Y float32 }{0.05, 0.70}
	loginButtonPos := struct{ X, Y float32 }{0.18, 0.92}

	envVars, err := readEnvFile(".env")
	if err != nil {
		fmt.Printf("Error reading .env file: %v\n", err)
		return
	}

	// username := envVars["USERNAME"]
	password := envVars["PASSWORD"]
	otp_secret := envVars["OTPSECRET"]

	totp, _ := generateTOTP(otp_secret)
	totp_text := strconv.Itoa(totp)

	hwnd, err := WaitForWindow(windowTitle, 30*time.Second)
	if err != nil {
		fmt.Println("Error finding window:", err)
		return
	}

	SetForegroundWindow(hwnd)
	time.Sleep(500 * time.Millisecond)

	rect, err := GetWindowRect(hwnd)
	if err != nil {
		fmt.Println("Error getting window rect:", err)
		return
	}

	windowWidth := rect.Right - rect.Left
	windowHeight := rect.Bottom - rect.Top

	// usernameX := rect.Left + int32(float32(windowWidth)*usernameFieldPos.X)
	// usernameY := rect.Top + int32(float32(windowHeight)*usernameFieldPos.Y)

	passwordX := rect.Left + int32(float32(windowWidth)*passwordFieldPos.X)
	passwordY := rect.Top + int32(float32(windowHeight)*passwordFieldPos.Y)

	otpX := rect.Left + int32(float32(windowWidth)*otpFieldPos.X)
	otpY := rect.Top + int32(float32(windowHeight)*otpFieldPos.Y)

	loginX := rect.Left + int32(float32(windowWidth)*loginButtonPos.X)
	loginY := rect.Top + int32(float32(windowHeight)*loginButtonPos.Y)

	fmt.Println("Starting automation...")

	// MouseClick(usernameX, usernameY)
	// TypeText(username)
	// time.Sleep(500 * time.Millisecond)

	MouseClick(passwordX, passwordY)
	TypeText(password)
	time.Sleep(500 * time.Millisecond)

	MouseClick(otpX, otpY)
	TypeText(totp_text)
	time.Sleep(500 * time.Millisecond)

	// Click login button
	MouseClick(loginX, loginY)

	fmt.Println("Automation completed")
}
