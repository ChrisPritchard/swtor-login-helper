package main

import (
	"fmt"
	"time"
	"unicode"
)

func SetCursorPos(x, y int32) bool {
	ret, _, _ := procSetCursorPos.Call(
		uintptr(x),
		uintptr(y),
	)
	return ret != 0
}

func MouseClick(x, y int32) {
	SetCursorPos(x, y)
	time.Sleep(100 * time.Millisecond)
	procMouseEvent.Call(
		uintptr(MOUSEEVENTF_LEFTDOWN),
		0, 0, 0, 0,
	)
	time.Sleep(50 * time.Millisecond)
	procMouseEvent.Call(
		uintptr(MOUSEEVENTF_LEFTUP),
		0, 0, 0, 0,
	)
	time.Sleep(100 * time.Millisecond)
}

var keyCodeMap = map[rune]uint16{
	'0': 0x30, '1': 0x31, '2': 0x32, '3': 0x33, '4': 0x34,
	'5': 0x35, '6': 0x36, '7': 0x37, '8': 0x38, '9': 0x39,

	'a': 0x41, 'b': 0x42, 'c': 0x43, 'd': 0x44, 'e': 0x45,
	'f': 0x46, 'g': 0x47, 'h': 0x48, 'i': 0x49, 'j': 0x4A,
	'k': 0x4B, 'l': 0x4C, 'm': 0x4D, 'n': 0x4E, 'o': 0x4F,
	'p': 0x50, 'q': 0x51, 'r': 0x52, 's': 0x53, 't': 0x54,
	'u': 0x55, 'v': 0x56, 'w': 0x57, 'x': 0x58, 'y': 0x59,
	'z': 0x5A,

	' ':  0x20,            // Space
	'\t': 0x09,            // Tab
	'\n': 0x0D,            // Enter
	'!':  0x31,            // Shift+1
	'@':  0x32,            // Shift+2
	'#':  0x33,            // Shift+3
	'$':  0x34,            // Shift+4
	'%':  0x35,            // Shift+5
	'^':  0x36,            // Shift+6
	'&':  0x37,            // Shift+7
	'*':  0x38,            // Shift+8
	'(':  0x39,            // Shift+9
	')':  0x30,            // Shift+0
	'-':  0xBD, '_': 0xBD, // Shift+-
	'=': 0xBB, '+': 0xBB, // Shift+=
	'[': 0xDB, '{': 0xDB, // Shift+[
	']': 0xDD, '}': 0xDD, // Shift+]
	'\\': 0xDC, '|': 0xDC, // Shift+\
	';': 0xBA, ':': 0xBA, // Shift+;
	'\'': 0xDE, '"': 0xDE, // Shift+'
	',': 0xBC, '<': 0xBC, // Shift+,
	'.': 0xBE, '>': 0xBE, // Shift+.
	'/': 0xBF, '?': 0xBF, // Shift+/
	'`': 0xC0, '~': 0xC0, // Shift+`
}

func SendKey(keyCode uint16, extended bool) {
	flags := uintptr(KEYEVENTF_KEYDOWN)
	if extended {
		flags |= KEYEVENTF_EXTENDEDKEY
	}

	procKeybdEvent.Call(
		uintptr(keyCode),
		0,
		flags,
		0,
	)
	time.Sleep(10 * time.Millisecond)

	procKeybdEvent.Call(
		uintptr(keyCode),
		0,
		uintptr(KEYEVENTF_KEYUP|flags),
		0,
	)
	time.Sleep(10 * time.Millisecond)
}

func PressShift() {
	procKeybdEvent.Call(
		uintptr(VK_SHIFT),
		0,
		uintptr(KEYEVENTF_KEYDOWN),
		0,
	)
	time.Sleep(10 * time.Millisecond)
}

func ReleaseShift() {
	procKeybdEvent.Call(
		uintptr(VK_SHIFT),
		0,
		uintptr(KEYEVENTF_KEYUP),
		0,
	)
	time.Sleep(10 * time.Millisecond)
}

func IsCapsLockOn() bool {
	ret, _, _ := procGetAsyncKeyState.Call(uintptr(VK_CAPITAL))
	return ret&0x0001 != 0
}

func TypeCharacter(c rune) {
	needsShift := unicode.IsUpper(c) ||
		(c >= '!' && c <= '&') ||
		(c >= '(' && c <= '+') ||
		(c >= ':' && c <= '>') ||
		(c == '?' || c == '~' || c == '{' || c == '}' || c == '|' || c == '"')

	capsOn := IsCapsLockOn()
	if unicode.IsLetter(c) {
		needsShift = (unicode.IsUpper(c) && !capsOn) || (unicode.IsLower(c) && capsOn)
	}

	lowerChar := unicode.ToLower(c)
	keyCode, ok := keyCodeMap[lowerChar]
	if !ok {
		fmt.Printf("Unsupported character: %c\n", c)
		return
	}

	if needsShift {
		PressShift()
	}

	SendKey(keyCode, false)

	if needsShift {
		ReleaseShift()
	}
}

func TypeText(text string) {
	for _, c := range text {
		TypeCharacter(c)
		time.Sleep(20 * time.Millisecond) // Small delay between characters
	}
}
