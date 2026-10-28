//go:build !windows

// 非 Windows：給開發時測試共用邏輯用的簡單命令列
// 用法：installer install <LIVE> [語言包]  /  installer uninstall <LIVE>  /  installer status <LIVE>
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("用法：installer install|uninstall|status <LIVE> [語言包]")
		os.Exit(2)
	}
	live, ok := normalizeGameFolder(os.Args[2])
	if !ok {
		fmt.Println("不是遊戲資料夾：", live)
		os.Exit(1)
	}
	log := func(s string) { fmt.Println(s) }
	var err error
	switch os.Args[1] {
	case "status":
		fmt.Println("g_language =", currentLanguage(live), "已安裝中文：", isChineseInstalled(live))
	case "install":
		var p *pack
		if p, err = findPack(os.Args[3:]); err == nil {
			fmt.Println("語言包：", p.Path, p.Version, p.Lines)
			err = install(live, p, log)
		}
	case "uninstall":
		err = uninstall(live, log)
	}
	if err != nil {
		fmt.Println("✘", err)
		os.Exit(1)
	}
}
