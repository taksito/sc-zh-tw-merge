package main

import "syscall"

// 讓命令提示字元正確顯示中文（UTF-8）
func setupConsole() {
	k := syscall.NewLazyDLL("kernel32.dll")
	k.NewProc("SetConsoleOutputCP").Call(65001)
	k.NewProc("SetConsoleCP").Call(65001)
}
