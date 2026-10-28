//go:build windows

// SC 繁中語言包安裝器：Windows 視窗介面（只用系統內建元件，不需額外安裝）
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

func init() { runtime.LockOSThread() } // 視窗訊息必須在同一個執行緒處理

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	comdlg32 = syscall.NewLazyDLL("comdlg32.dll")
	ole32    = syscall.NewLazyDLL("ole32.dll")
	comctl32 = syscall.NewLazyDLL("comctl32.dll")

	pRegisterClassExW      = user32.NewProc("RegisterClassExW")
	pCreateWindowExW       = user32.NewProc("CreateWindowExW")
	pDefWindowProcW        = user32.NewProc("DefWindowProcW")
	pGetMessageW           = user32.NewProc("GetMessageW")
	pTranslateMessage      = user32.NewProc("TranslateMessage")
	pDispatchMessageW      = user32.NewProc("DispatchMessageW")
	pIsDialogMessageW      = user32.NewProc("IsDialogMessageW")
	pPostQuitMessage       = user32.NewProc("PostQuitMessage")
	pSendMessageW          = user32.NewProc("SendMessageW")
	pSetWindowTextW        = user32.NewProc("SetWindowTextW")
	pMessageBoxW           = user32.NewProc("MessageBoxW")
	pEnableWindow          = user32.NewProc("EnableWindow")
	pShowWindow            = user32.NewProc("ShowWindow")
	pUpdateWindow          = user32.NewProc("UpdateWindow")
	pLoadCursorW           = user32.NewProc("LoadCursorW")
	pLoadIconW             = user32.NewProc("LoadIconW")
	pAdjustWindowRectEx    = user32.NewProc("AdjustWindowRectEx")
	pGetDpiForSystem       = user32.NewProc("GetDpiForSystem")
	pSetProcessDpiAwareCtx = user32.NewProc("SetProcessDpiAwarenessContext")
	pSetProcessDPIAware    = user32.NewProc("SetProcessDPIAware")
	pGetSystemMetrics      = user32.NewProc("GetSystemMetrics")
	pCreateFontW           = gdi32.NewProc("CreateFontW")
	pGetModuleHandleW      = kernel32.NewProc("GetModuleHandleW")
	pCreateActCtxW         = kernel32.NewProc("CreateActCtxW")
	pActivateActCtx        = kernel32.NewProc("ActivateActCtx")
	pSHBrowseForFolderW    = shell32.NewProc("SHBrowseForFolderW")
	pSHGetPathFromIDListW  = shell32.NewProc("SHGetPathFromIDListW")
	pDragAcceptFiles       = shell32.NewProc("DragAcceptFiles")
	pDragQueryFileW        = shell32.NewProc("DragQueryFileW")
	pDragFinish            = shell32.NewProc("DragFinish")
	pShellExecuteW         = shell32.NewProc("ShellExecuteW")
	pGetOpenFileNameW      = comdlg32.NewProc("GetOpenFileNameW")
	pCoInitializeEx        = ole32.NewProc("CoInitializeEx")
	pCoTaskMemFree         = ole32.NewProc("CoTaskMemFree")
	pInitCommonControlsEx  = comctl32.NewProc("InitCommonControlsEx")
)

const (
	wmCreate    = 0x0001
	wmDestroy   = 0x0002
	wmSetFont   = 0x0030
	wmCommand   = 0x0111
	wmDropFiles = 0x0233

	wsOverlapped = 0x00000000
	wsCaption    = 0x00C00000
	wsSysMenu    = 0x00080000
	wsMinimize   = 0x00020000
	wsChild      = 0x40000000
	wsVisible    = 0x10000000
	wsTabStop    = 0x00010000
	wsVScroll    = 0x00200000
	wsExClient   = 0x00000200
	wsExAccept   = 0x00000010

	cbsDropDownList = 0x0003
	cbAddString     = 0x0143
	cbResetContent  = 0x014B
	cbSetCurSel     = 0x014E
	cbGetCurSel     = 0x0147
	cbnSelChange    = 1

	esMultiline   = 0x0004
	esAutoVScroll = 0x0040
	esAutoHScroll = 0x0080
	esReadOnly    = 0x0800
	emSetSel      = 0x00B1
	emReplaceSel  = 0x00C2

	bsDefPush    = 0x0001
	ssPathEllips = 0x8000

	mbOK       = 0x00
	mbYesNo    = 0x04
	mbIconErr  = 0x10
	mbIconQ    = 0x20
	mbIconWarn = 0x30
	mbIconInfo = 0x40
	idYes      = 6
)

const (
	idCombo = 101 + iota
	idBrowseGame
	idPackEdit
	idBrowsePack
	idInstall
	idUninstall
	idOpenBackup
	idLog
	idStatus
	idPackInfo
)

type wndClassEx struct {
	cbSize        uint32
	style         uint32
	lpfnWndProc   uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     uintptr
	hIcon         uintptr
	hCursor       uintptr
	hbrBackground uintptr
	lpszMenuName  *uint16
	lpszClassName *uint16
	hIconSm       uintptr
}

type point struct{ x, y int32 }
type msgT struct {
	hwnd    uintptr
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      point
	private uint32
}
type rect struct{ left, top, right, bottom int32 }

type browseInfo struct {
	hwndOwner      uintptr
	pidlRoot       uintptr
	pszDisplayName *uint16
	lpszTitle      *uint16
	ulFlags        uint32
	lpfn           uintptr
	lParam         uintptr
	iImage         int32
}

type openFileName struct {
	lStructSize       uint32
	hwndOwner         uintptr
	hInstance         uintptr
	lpstrFilter       *uint16
	lpstrCustomFilter *uint16
	nMaxCustFilter    uint32
	nFilterIndex      uint32
	lpstrFile         *uint16
	nMaxFile          uint32
	lpstrFileTitle    *uint16
	nMaxFileTitle     uint32
	lpstrInitialDir   *uint16
	lpstrTitle        *uint16
	flags             uint32
	nFileOffset       uint16
	nFileExtension    uint16
	lpstrDefExt       *uint16
	lCustData         uintptr
	lpfnHook          uintptr
	lpTemplateName    *uint16
	pvReserved        uintptr
	dwReserved        uint32
	flagsEx           uint32
}

type actCtx struct {
	cbSize                 uint32
	dwFlags                uint32
	lpSource               *uint16
	wProcessorArchitecture uint16
	wLangId                uint16
	lpAssemblyDirectory    *uint16
	lpResourceName         *uint16
	lpApplicationName      *uint16
	hModule                uintptr
}

type initCC struct{ size, icc uint32 }

func u16(s string) *uint16 { p, _ := syscall.UTF16PtrFromString(s); return p }

// ---------- 狀態 ----------

var (
	hInst, hMain, hFont                                               uintptr
	hCombo, hPackEdit, hStatus, hPackInfo, hLog, hInstall, hUninstall uintptr
	dpi                                                               uintptr = 96
	folders                                                           []string
	curPack                                                           *pack
)

func s(x int) uintptr { return uintptr(x) * dpi / 96 }

func send(h uintptr, m uint32, w, l uintptr) uintptr {
	r, _, _ := pSendMessageW.Call(h, uintptr(m), w, l)
	return r
}
func setText(h uintptr, t string) { pSetWindowTextW.Call(h, uintptr(unsafe.Pointer(u16(t)))) }
func enable(h uintptr, on bool) {
	v := uintptr(0)
	if on {
		v = 1
	}
	pEnableWindow.Call(h, v)
}
func msgBox(text string, flags uintptr) uintptr {
	r, _, _ := pMessageBoxW.Call(hMain, uintptr(unsafe.Pointer(u16(text))), uintptr(unsafe.Pointer(u16(appName))), flags)
	return r
}
func logLine(t string) {
	t = strings.ReplaceAll(t, "\n", "\r\n") + "\r\n"
	send(hLog, emSetSel, ^uintptr(0), ^uintptr(0))
	send(hLog, emReplaceSel, 0, uintptr(unsafe.Pointer(u16(t))))
}

func control(class, text string, style uintptr, exStyle uintptr, x, y, w, h int, id int) uintptr {
	hwnd, _, _ := pCreateWindowExW.Call(exStyle, uintptr(unsafe.Pointer(u16(class))), uintptr(unsafe.Pointer(u16(text))),
		wsChild|wsVisible|style, s(x), s(y), s(w), s(h), hMain, uintptr(id), hInst, 0)
	send(hwnd, wmSetFont, hFont, 1)
	return hwnd
}

// ---------- 畫面更新 ----------

func selectedFolder() string {
	i := int(send(hCombo, cbGetCurSel, 0, 0))
	if i < 0 || i >= len(folders) {
		return ""
	}
	return folders[i]
}

func fillFolders(selectPath string) {
	send(hCombo, cbResetContent, 0, 0)
	sel := 0
	for i, f := range folders {
		send(hCombo, cbAddString, 0, uintptr(unsafe.Pointer(u16(f))))
		if strings.EqualFold(f, selectPath) {
			sel = i
		}
	}
	if len(folders) > 0 {
		send(hCombo, cbSetCurSel, uintptr(sel), 0)
	}
	refreshStatus()
}

func refreshStatus() {
	live := selectedFolder()
	switch {
	case live == "":
		setText(hStatus, "目前狀態：沒有找到遊戲資料夾，請按右邊「瀏覽…」選擇 StarCitizen\\LIVE")
		setText(hInstall, "安裝中文化")
		enable(hInstall, false)
		enable(hUninstall, false)
		return
	case isChineseInstalled(live):
		setText(hStatus, "目前狀態：已安裝繁體中文")
		setText(hInstall, "更新中文化")
		enable(hUninstall, true)
	default:
		cur := currentLanguage(live)
		if cur == "" {
			cur = "english"
		}
		setText(hStatus, "目前狀態：英文（g_language = "+cur+"）")
		setText(hInstall, "安裝中文化")
		enable(hUninstall, false)
	}
	enable(hInstall, curPack != nil)
}

func setPack(p *pack) {
	curPack = p
	if p == nil {
		setText(hPackEdit, "")
		setText(hPackInfo, "找不到語言包：請從網站下載 global_zh-TW.zip，或按「瀏覽…」選擇，也可以直接把檔案拖進這個視窗")
	} else {
		setText(hPackEdit, p.Path)
		setText(hPackInfo, fmt.Sprintf("%s，共 %d 條", p.Version, p.Lines))
	}
	refreshStatus()
}

func usePackPath(path string) {
	p, err := loadPack(path)
	if err != nil {
		msgBox("無法使用這個檔案：\n"+path+"\n\n"+err.Error(), mbOK|mbIconWarn)
		return
	}
	setPack(p)
	logLine("已選擇語言包：" + p.Path)
}

func useFolderPath(path string) {
	live, ok := normalizeGameFolder(path)
	if !ok {
		msgBox("這個資料夾裡找不到 Data.p4k 或 Bin64\\StarCitizen.exe：\n"+path+"\n\n請選擇 StarCitizen 底下的 LIVE 資料夾。", mbOK|mbIconWarn)
		return
	}
	found := false
	for _, f := range folders {
		if strings.EqualFold(f, live) {
			found = true
		}
	}
	if !found {
		folders = append(folders, live)
	}
	fillFolders(live)
}

// ---------- 動作 ----------

func doInstall() {
	live := selectedFolder()
	if live == "" || curPack == nil {
		return
	}
	verb := "安裝"
	if isChineseInstalled(live) {
		verb = "更新"
	}
	q := fmt.Sprintf("要%s繁體中文語言包嗎？\n\n遊戲資料夾：%s\n語言包：%s（%s）\n\n舊檔案會先備份到「%s」。\n請先關閉遊戲。", verb, live, filepath.Base(curPack.Path), curPack.Version, backupDir)
	if msgBox(q, mbYesNo|mbIconQ) != idYes {
		return
	}
	logLine("── " + verb + "中文化 ──")
	if err := install(live, curPack, logLine); err != nil {
		logLine("✘ " + err.Error())
		offerElevate(err)
	} else {
		msgBox(verb+"完成！啟動遊戲就會是繁體中文。", mbOK|mbIconInfo)
	}
	refreshStatus()
}

func doUninstall() {
	live := selectedFolder()
	if live == "" {
		return
	}
	q := "要移除中文化、變回英文嗎？\n\n只會把 user.cfg 的 g_language 改成 english（先備份），\n中文語言包會保留，之後可隨時再安裝回來。\n請先關閉遊戲。"
	if msgBox(q, mbYesNo|mbIconQ) != idYes {
		return
	}
	logLine("── 移除中文化 ──")
	if err := uninstall(live, logLine); err != nil {
		logLine("✘ " + err.Error())
		offerElevate(err)
	} else {
		msgBox("已變回英文。", mbOK|mbIconInfo)
	}
	refreshStatus()
}

// 權限不足時，問要不要以系統管理員身分重新開啟
func offerElevate(err error) {
	if !strings.Contains(err.Error(), "Access is denied") && !strings.Contains(strings.ToLower(err.Error()), "permission") && !strings.Contains(err.Error(), "拒絕存取") {
		msgBox(err.Error(), mbOK|mbIconErr)
		return
	}
	if msgBox(err.Error()+"\n\n要以系統管理員身分重新開啟安裝器嗎？", mbYesNo|mbIconWarn) != idYes {
		return
	}
	exe, e := os.Executable()
	if e != nil {
		return
	}
	args := ""
	if curPack != nil {
		args = `"` + curPack.Path + `"`
	}
	r, _, _ := pShellExecuteW.Call(hMain, uintptr(unsafe.Pointer(u16("runas"))), uintptr(unsafe.Pointer(u16(exe))), uintptr(unsafe.Pointer(u16(args))), 0, 1)
	if r > 32 { // 成功開啟新的（管理員）視窗，就關掉目前這個
		pPostQuitMessage.Call(0)
	}
}

func doOpenBackup() {
	live := selectedFolder()
	if live == "" {
		return
	}
	dir := backupRoot(live)
	if _, err := os.Stat(dir); err != nil {
		msgBox("還沒有任何備份。\n安裝或移除時才會建立：\n"+dir, mbOK|mbIconInfo)
		return
	}
	pShellExecuteW.Call(hMain, uintptr(unsafe.Pointer(u16("open"))), uintptr(unsafe.Pointer(u16(dir))), 0, 0, 1)
}

func browseFolder() {
	buf := make([]uint16, 260)
	bi := browseInfo{hwndOwner: hMain, pszDisplayName: &buf[0], lpszTitle: u16("選擇 Star Citizen 的 LIVE 資料夾（或 StarCitizen 資料夾）"), ulFlags: 0x0001 | 0x0040}
	pidl, _, _ := pSHBrowseForFolderW.Call(uintptr(unsafe.Pointer(&bi)))
	if pidl == 0 {
		return
	}
	path := make([]uint16, 1024)
	pSHGetPathFromIDListW.Call(pidl, uintptr(unsafe.Pointer(&path[0])))
	pCoTaskMemFree.Call(pidl)
	useFolderPath(syscall.UTF16ToString(path))
}

func browsePack() {
	filter := utf16z("語言包 (*.zip; *.ini; *.txt)", "*.zip;*.ini;*.txt", "所有檔案", "*.*")
	file := make([]uint16, 1024)
	init := ""
	if home, err := os.UserHomeDir(); err == nil {
		init = filepath.Join(home, "Downloads")
	}
	ofn := openFileName{hwndOwner: hMain, lpstrFilter: &filter[0], nFilterIndex: 1, lpstrFile: &file[0], nMaxFile: uint32(len(file)),
		lpstrInitialDir: u16(init), lpstrTitle: u16("選擇語言包（從網站下載的 zip 或 global.ini）"), flags: 0x1000 | 0x0800 | 0x00080000}
	ofn.lStructSize = uint32(unsafe.Sizeof(ofn))
	if r, _, _ := pGetOpenFileNameW.Call(uintptr(unsafe.Pointer(&ofn))); r != 0 {
		usePackPath(syscall.UTF16ToString(file))
	}
}

// 多個字串以 \0 分隔、最後再加一個 \0
func utf16z(parts ...string) []uint16 {
	var out []uint16
	for _, p := range parts {
		out = append(out, syscall.StringToUTF16(p)...)
	}
	return append(out, 0)
}

func onDrop(hDrop uintptr) {
	defer pDragFinish.Call(hDrop)
	n, _, _ := pDragQueryFileW.Call(hDrop, 0xFFFFFFFF, 0, 0)
	for i := uintptr(0); i < n; i++ {
		buf := make([]uint16, 1024)
		pDragQueryFileW.Call(hDrop, i, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		p := syscall.UTF16ToString(buf)
		if st, err := os.Stat(p); err == nil && st.IsDir() {
			useFolderPath(p)
		} else {
			usePackPath(p)
		}
	}
}

// ---------- 視窗 ----------

const clientW, clientH = 600, 486

func buildUI() {
	control("STATIC", "遊戲資料夾", 0, 0, 16, 14, 568, 20, 0)
	hCombo = control("COMBOBOX", "", cbsDropDownList|wsTabStop|wsVScroll, 0, 16, 36, 452, 300, idCombo)
	control("BUTTON", "瀏覽…", wsTabStop, 0, 476, 35, 108, 28, idBrowseGame)
	hStatus = control("STATIC", "", 0, 0, 16, 70, 568, 20, idStatus)

	control("STATIC", "語言包", 0, 0, 16, 104, 568, 20, 0)
	hPackEdit = control("EDIT", "", esReadOnly|esAutoHScroll|wsTabStop, wsExClient, 16, 126, 452, 26, idPackEdit)
	control("BUTTON", "瀏覽…", wsTabStop, 0, 476, 125, 108, 28, idBrowsePack)
	hPackInfo = control("STATIC", "", ssPathEllips, 0, 16, 158, 568, 20, idPackInfo)

	hInstall = control("BUTTON", "安裝中文化", bsDefPush|wsTabStop, 0, 16, 194, 190, 38, idInstall)
	hUninstall = control("BUTTON", "移除中文化（變回英文）", wsTabStop, 0, 214, 194, 190, 38, idUninstall)
	control("BUTTON", "開啟備份資料夾", wsTabStop, 0, 412, 194, 172, 38, idOpenBackup)

	control("STATIC", "訊息", 0, 0, 16, 246, 568, 20, 0)
	hLog = control("EDIT", "", esMultiline|esAutoVScroll|esReadOnly|wsVScroll|wsTabStop, wsExClient, 16, 268, 568, 170, idLog)
	control("STATIC", "安裝前請先關閉遊戲。可以把下載的 zip 或遊戲資料夾直接拖進這個視窗。", 0, 0, 16, 448, 568, 20, 0)

	logLine(appName + " v" + version)
	logLine("正在尋找 Star Citizen…")
	folders = findGameFolders()
	if len(folders) == 0 {
		logLine("沒有自動找到遊戲資料夾，請按「瀏覽…」選擇 StarCitizen\\LIVE。")
	} else {
		logLine(fmt.Sprintf("找到 %d 個遊戲資料夾。", len(folders)))
	}
	p, err := findPack(os.Args[1:])
	if err == nil {
		logLine("語言包：" + p.Path)
	}
	fillFolders("")
	setPack(p)
}

func wndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	switch msg {
	case wmCreate:
		hMain = hwnd
		buildUI()
		return 0
	case wmCommand:
		id, code := wParam&0xFFFF, (wParam>>16)&0xFFFF
		switch id {
		case idCombo:
			if code == cbnSelChange {
				refreshStatus()
			}
		case idBrowseGame:
			browseFolder()
		case idBrowsePack:
			browsePack()
		case idInstall:
			doInstall()
		case idUninstall:
			doUninstall()
		case idOpenBackup:
			doOpenBackup()
		}
		return 0
	case wmDropFiles:
		onDrop(wParam)
		return 0
	case wmDestroy:
		pPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := pDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

// 用執行期的 manifest 啟用新版控制項外觀（不需要把資源編進 exe）
func enableVisualStyles() {
	const manifest = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <dependency><dependentAssembly>
    <assemblyIdentity type="win32" name="Microsoft.Windows.Common-Controls" version="6.0.0.0" processorArchitecture="*" publicKeyToken="6595b64144ccf1df" language="*"/>
  </dependentAssembly></dependency>
</assembly>`
	path := filepath.Join(os.TempDir(), "sc-zhtw-installer.manifest")
	if os.WriteFile(path, []byte(manifest), 0o644) != nil {
		return
	}
	ac := actCtx{lpSource: u16(path)}
	ac.cbSize = uint32(unsafe.Sizeof(ac))
	h, _, _ := pCreateActCtxW.Call(uintptr(unsafe.Pointer(&ac)))
	if h == 0 || h == ^uintptr(0) {
		return
	}
	var cookie uintptr
	pActivateActCtx.Call(h, uintptr(unsafe.Pointer(&cookie)))
	icc := initCC{size: 8, icc: 0x4000 | 0xFF} // ICC_STANDARD_CLASSES 等
	pInitCommonControlsEx.Call(uintptr(unsafe.Pointer(&icc)))
}

func main() {
	// 高解析度螢幕不要模糊
	if pSetProcessDpiAwareCtx.Find() == nil {
		pSetProcessDpiAwareCtx.Call(^uintptr(1)) // DPI_AWARENESS_CONTEXT_SYSTEM_AWARE (-2)
	} else if pSetProcessDPIAware.Find() == nil {
		pSetProcessDPIAware.Call()
	}
	if pGetDpiForSystem.Find() == nil {
		if d, _, _ := pGetDpiForSystem.Call(); d >= 96 {
			dpi = d
		}
	}
	enableVisualStyles()
	pCoInitializeEx.Call(0, 2) // 資料夾選擇視窗需要

	hInst, _, _ = pGetModuleHandleW.Call(0)
	height := -int64((10*int(dpi) + 36) / 72) // 約 10pt；負值代表字元高度
	hFont, _, _ = pCreateFontW.Call(uintptr(height), 0, 0, 0, 400, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(u16("Microsoft JhengHei UI"))))

	cursor, _, _ := pLoadCursorW.Call(0, 32512)
	icon, _, _ := pLoadIconW.Call(0, 32512)
	cls := u16("SCzhTWInstaller")
	wc := wndClassEx{lpfnWndProc: syscall.NewCallback(wndProc), hInstance: hInst, hIcon: icon, hIconSm: icon, hCursor: cursor, hbrBackground: 15 + 1, lpszClassName: cls}
	wc.cbSize = uint32(unsafe.Sizeof(wc))
	pRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))

	style := uintptr(wsOverlapped | wsCaption | wsSysMenu | wsMinimize)
	r := rect{0, 0, int32(s(clientW)), int32(s(clientH))}
	pAdjustWindowRectEx.Call(uintptr(unsafe.Pointer(&r)), style, 0, wsExAccept)
	w, h := uintptr(r.right-r.left), uintptr(r.bottom-r.top)
	sw, _, _ := pGetSystemMetrics.Call(0)
	sh, _, _ := pGetSystemMetrics.Call(1)
	x, y := (sw-w)/2, (sh-h)/2
	if int(sw) < int(w) {
		x = 0
	}
	if int(sh) < int(h) {
		y = 0
	}
	hwnd, _, _ := pCreateWindowExW.Call(wsExAccept, uintptr(unsafe.Pointer(cls)), uintptr(unsafe.Pointer(u16(appName+" v"+version))),
		style, x, y, w, h, 0, 0, hInst, 0)
	if hwnd == 0 {
		return
	}
	pDragAcceptFiles.Call(hwnd, 1)
	pShowWindow.Call(hwnd, 1)
	pUpdateWindow.Call(hwnd)

	var m msgT
	for {
		ret, _, _ := pGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(ret) <= 0 {
			break
		}
		if r, _, _ := pIsDialogMessageW.Call(hwnd, uintptr(unsafe.Pointer(&m))); r != 0 {
			continue // 讓 Tab 鍵可以切換按鈕
		}
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		pDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
}
