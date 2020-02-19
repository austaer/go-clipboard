package main

import (
	"flag"
	"syscall"
	"unsafe"
)

const (
	GHND           = 0x0042
	GMEM_FIXED     = 0x0000
	GMEM_MOVEABLE  = 0x0002
	GMEM_ZEROINIT  = 0x0040
	GPTR           = GMEM_FIXED | GMEM_ZEROINIT
	CF_UNICODETEXT = 0x000D
)

var (
	user32                     = syscall.MustLoadDLL("user32.dll")
	openClipboard              = user32.MustFindProc("OpenClipboard")
	closeClipboard             = user32.MustFindProc("CloseClipboard")
	emptyClipboard             = user32.MustFindProc("EmptyClipboard")
	getClipboardData           = user32.MustFindProc("GetClipboardData")
	setClipboardData           = user32.MustFindProc("SetClipboardData")
	isClipboardFormatAvailable = user32.MustFindProc("IsClipboardFormatAvailable")

	kernel32     = syscall.NewLazyDLL("kernel32")
	globalAlloc  = kernel32.NewProc("GlobalAlloc")
	globalFree   = kernel32.NewProc("GlobalFree")
	globalLock   = kernel32.NewProc("GlobalLock")
	globalUnlock = kernel32.NewProc("GlobalUnlock")
	lstrcpy      = kernel32.NewProc("lstrcpyW")
	copyMemory   = kernel32.NewProc("CopyMemory")
)

func main() {
	input := flag.String("p", "", "輸入要被複製的值")
	flag.Parse()
	value := *input
	setClipBoard(value)
}

func setClipBoard(str string) {
	r, _, err := openClipboard.Call(0)
	if r == 0 {
		alert("openClipboard error ", "error")
		panic(err)
	}
	defer closeClipboard.Call()
	emptyClipboard.Call()

	iLen := len([]byte(str))*2 + 2
	iStrPtr, _, _ := globalAlloc.Call(GPTR, IntPtr(iLen))
	iLock, _, _ := globalLock.Call(iStrPtr)
	lstrcpy.Call(iLock, StrPtr(str))
	globalUnlock.Call(iStrPtr)
	setClipboardData.Call(CF_UNICODETEXT, iStrPtr)
}

func alert(text string, title string) {
	inUser32 := syscall.NewLazyDLL("user32.dll")
	MessageBoxW := inUser32.NewProc("MessageBoxW")
	MessageBoxW.Call(IntPtr(0), StrPtr(text), StrPtr(title), IntPtr(0))
}

func IntPtr(n int) uintptr {
	return uintptr(n)
}

func StrPtr(s string) uintptr {
	return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s)))
}
