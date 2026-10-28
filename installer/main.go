// SC 繁中語言包安裝器
// 找到 Star Citizen 的 LIVE 資料夾，把 global.ini 放進
// Data\Localization\chinese_(traditional)\，並設定 user.cfg。
// 舊檔案會先備份到 LIVE\語言包備份\日期時間\。
package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	langFolder = "chinese_(traditional)"
	backupDir  = "語言包備份"
	version    = "1.1.0"
)

var in = bufio.NewReader(os.Stdin)

func main() {
	setupConsole()
	defer pause()

	fmt.Println("==============================================")
	fmt.Println("  Star Citizen 繁體中文語言包安裝器  v" + version)
	fmt.Println("==============================================")
	fmt.Println()

	// 1. 找遊戲資料夾
	live := pickGameFolder()
	if live == "" {
		return
	}
	cfg := filepath.Join(live, "user.cfg")
	target := filepath.Join(live, "Data", "Localization", langFolder, "global.ini")

	// 2. 已經裝了中文 → 讓使用者選更新或移除
	if cur := currentLanguage(cfg); strings.EqualFold(cur, langFolder) {
		fmt.Println()
		fmt.Println("偵測到這個資料夾已安裝繁體中文語言包。")
		fmt.Println("  [1] 更新中文化")
		fmt.Println("  [2] 移除中文化（變回英文）")
		fmt.Println("  [0] 取消")
		for {
			switch prompt("請選擇：") {
			case "1":
				install(live, cfg, target)
				return
			case "2":
				uninstall(live, cfg)
				return
			case "0", "":
				fmt.Println("已取消，沒有修改任何檔案。")
				return
			default:
				fmt.Println("請輸入 1、2 或 0。")
			}
		}
	}
	install(live, cfg, target)
}

// 安裝／更新中文化
func install(live, cfg, target string) {
	fmt.Println()
	src, data, err := findPackage()
	if err != nil {
		fail(err.Error())
		return
	}
	fmt.Printf("語言包：%s\n", src)
	fmt.Printf("　　　　%s，%d 條\n", versionLine(data), countLines(data))

	fmt.Println()
	fmt.Println("即將進行：")
	fmt.Printf("  1. 寫入 %s\n", target)
	fmt.Printf("  2. 在 %s 設定 g_language = %s\n", cfg, langFolder)
	fmt.Printf("  3. 舊檔案先備份到 %s\n", filepath.Join(live, backupDir))
	fmt.Println()
	fmt.Println("請先關閉遊戲。")
	if !ask("確定要安裝嗎？(Y/N)：") {
		fmt.Println("已取消，沒有修改任何檔案。")
		return
	}
	if !backup(live, target, cfg) {
		return
	}

	// 寫入 global.ini（確保 UTF-8 with BOM）
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		fail("建立資料夾失敗：" + err.Error() + "\n（可能需要用「以系統管理員身分執行」）")
		return
	}
	if !bytes.HasPrefix(data, []byte{0xEF, 0xBB, 0xBF}) {
		data = append([]byte{0xEF, 0xBB, 0xBF}, data...)
	}
	if err := os.WriteFile(target, data, 0o644); err != nil {
		fail("寫入 global.ini 失敗：" + err.Error() + "\n（可能需要用「以系統管理員身分執行」）")
		return
	}
	fmt.Println("✔ 已寫入 global.ini")

	if err := setLanguage(cfg, langFolder, true); err != nil {
		fail("修改 user.cfg 失敗：" + err.Error())
		return
	}
	fmt.Println("✔ 已設定 user.cfg")
	fmt.Println()
	fmt.Println("完成！啟動遊戲就會是繁體中文。")
	fmt.Println("要還原：把「" + backupDir + "」裡對應日期的檔案複製回 LIVE 資料夾即可。")
}

// 移除中文化：只把 user.cfg 的 g_language 改回 english
func uninstall(live, cfg string) {
	fmt.Println()
	fmt.Println("即將進行：")
	fmt.Printf("  1. 在 %s 把 g_language 改成 english\n", cfg)
	fmt.Printf("  2. 舊的 user.cfg 先備份到 %s\n", filepath.Join(live, backupDir))
	fmt.Println("  （中文 global.ini 會留著，之後想換回中文，再執行一次選「更新中文化」即可）")
	fmt.Println()
	fmt.Println("請先關閉遊戲。")
	if !ask("確定要變回英文嗎？(Y/N)：") {
		fmt.Println("已取消，沒有修改任何檔案。")
		return
	}
	if !backup(live, cfg) {
		return
	}
	if err := setLanguage(cfg, "english", false); err != nil {
		fail("修改 user.cfg 失敗：" + err.Error())
		return
	}
	fmt.Println("✔ 已把 user.cfg 改成英文")
	fmt.Println()
	fmt.Println("完成！啟動遊戲就會是英文。")
}

// 把存在的檔案備份到 LIVE\語言包備份\日期時間\
func backup(live string, files ...string) bool {
	bdir := filepath.Join(live, backupDir, time.Now().Format("20060102-150405"))
	for i := 2; ; i++ { // 同一秒執行兩次時不要蓋掉前一份備份
		if _, err := os.Stat(bdir); err != nil {
			break
		}
		bdir = filepath.Join(live, backupDir, time.Now().Format("20060102-150405")+"-"+strconv.Itoa(i))
	}
	n := 0
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			continue
		}
		rel, _ := filepath.Rel(live, f)
		if err := copyFile(f, filepath.Join(bdir, rel)); err != nil {
			fail("備份失敗：" + err.Error() + "\n（可能需要用「以系統管理員身分執行」）")
			return false
		}
		n++
	}
	if n > 0 {
		fmt.Printf("✔ 已備份 %d 個舊檔案到 %s\n", n, bdir)
	} else {
		fmt.Println("✔ 沒有舊檔案需要備份")
	}
	return true
}

// ---------- 語言包檔案 ----------

// 依序找：拖曳到程式上的檔案 → 程式旁邊 → 下載資料夾；取最新的 zip / global.ini / global.ini.txt
func findPackage() (string, []byte, error) {
	var cands []string
	for _, a := range os.Args[1:] {
		cands = append(cands, a)
	}
	if len(cands) == 0 {
		dirs := []string{exeDir()}
		if home, err := os.UserHomeDir(); err == nil {
			dirs = append(dirs, filepath.Join(home, "Downloads"))
		}
		type fi struct {
			p string
			t time.Time
		}
		var found []fi
		re := regexp.MustCompile(`(?i)^(global.*\.zip|global\.ini(\.txt)?|global \(\d+\)\.ini(\.txt)?|global\.ini \(\d+\)\.txt)$`)
		for _, d := range dirs {
			ents, _ := os.ReadDir(d)
			for _, e := range ents {
				if e.IsDir() || !re.MatchString(e.Name()) {
					continue
				}
				info, err := e.Info()
				if err == nil {
					found = append(found, fi{filepath.Join(d, e.Name()), info.ModTime()})
				}
			}
		}
		sort.Slice(found, func(i, j int) bool { return found[i].t.After(found[j].t) })
		for _, f := range found {
			cands = append(cands, f.p)
		}
	}
	for _, c := range cands {
		data, err := readPackage(c)
		if err == nil {
			return c, data, nil
		}
	}
	return "", nil, fmt.Errorf("找不到語言包。\n請把從網站下載的 global_zh-TW.zip（或 global.ini）\n放在這個程式旁邊或「下載」資料夾，或直接拖曳到程式圖示上。")
}

func readPackage(p string) ([]byte, error) {
	if strings.EqualFold(filepath.Ext(p), ".zip") {
		zr, err := zip.OpenReader(p)
		if err != nil {
			return nil, err
		}
		defer zr.Close()
		for _, f := range zr.File {
			if strings.EqualFold(filepath.Base(f.Name), "global.ini") {
				rc, err := f.Open()
				if err != nil {
					return nil, err
				}
				defer rc.Close()
				b, err := io.ReadAll(rc)
				if err != nil {
					return nil, err
				}
				return checkIni(b)
			}
		}
		return nil, fmt.Errorf("zip 裡沒有 global.ini")
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	return checkIni(b)
}

func checkIni(b []byte) ([]byte, error) {
	if len(b) < 100000 || !bytes.Contains(b[:min(len(b), 4096)], []byte("=")) {
		return nil, fmt.Errorf("不像是語言包檔案")
	}
	return b, nil
}

func versionLine(b []byte) string {
	m := regexp.MustCompile(`(?m)^Frontend_PU_Version=(.*)$`).FindSubmatch(b)
	if m == nil {
		return "版本不明"
	}
	return strings.TrimSpace(string(m[1]))
}

func countLines(b []byte) int { return bytes.Count(b, []byte("\n")) + 1 }

// ---------- 遊戲資料夾 ----------

func isGameFolder(p string) bool {
	for _, f := range []string{"Data.p4k", filepath.Join("Bin64", "StarCitizen.exe")} {
		if _, err := os.Stat(filepath.Join(p, f)); err == nil {
			return true
		}
	}
	return false
}

func findGameFolders() []string {
	seen := map[string]bool{}
	var out []string
	add := func(p string) {
		p = filepath.Clean(p)
		k := strings.ToLower(p)
		if seen[k] || !isGameFolder(p) {
			return
		}
		seen[k] = true
		out = append(out, p)
	}
	// RSI 啟動器的紀錄裡會有安裝路徑
	if ad := os.Getenv("APPDATA"); ad != "" {
		re := regexp.MustCompile(`(?i)[A-Z]:[\\/][^"'\r\n<>|*?]*?StarCitizen[\\/](LIVE|PTU|EPTU|TECH-PREVIEW|HOTFIX)`)
		for _, lf := range []string{
			filepath.Join(ad, "rsilauncher", "logs", "log.log"),
			filepath.Join(ad, "rsilauncher", "log.log"),
			filepath.Join(ad, "rsilauncher", "settings.json"),
		} {
			if b, err := os.ReadFile(lf); err == nil {
				for _, m := range re.FindAll(b, -1) {
					add(strings.ReplaceAll(strings.ReplaceAll(string(m), `\\`, `\`), "/", `\`))
				}
			}
		}
	}
	// 常見安裝位置 × 每個磁碟
	subs := []string{
		`Program Files\Roberts Space Industries\StarCitizen`,
		`Roberts Space Industries\StarCitizen`,
		`Games\Roberts Space Industries\StarCitizen`,
		`StarCitizen`,
		`Games\StarCitizen`,
	}
	for c := 'C'; c <= 'Z'; c++ {
		root := string(c) + `:\`
		if _, err := os.Stat(root); err != nil {
			continue
		}
		for _, s := range subs {
			for _, ch := range []string{"LIVE", "PTU", "EPTU", "HOTFIX"} {
				add(filepath.Join(root, s, ch))
			}
		}
	}
	// LIVE 排前面
	sort.SliceStable(out, func(i, j int) bool {
		return strings.HasSuffix(strings.ToUpper(out[i]), `\LIVE`) && !strings.HasSuffix(strings.ToUpper(out[j]), `\LIVE`)
	})
	return out
}

func pickGameFolder() string {
	fmt.Println("正在尋找 Star Citizen 安裝位置…")
	list := findGameFolders()
	for {
		if len(list) > 0 {
			fmt.Println("找到以下遊戲資料夾：")
			for i, p := range list {
				fmt.Printf("  [%d] %s\n", i+1, p)
			}
			fmt.Println("  [0] 手動輸入路徑")
			s := prompt(fmt.Sprintf("請選擇（直接按 Enter = %d）：", 1))
			if s == "" {
				s = "1"
			}
			n, err := strconv.Atoi(s)
			if err == nil && n >= 1 && n <= len(list) {
				return list[n-1]
			}
			if err != nil || n != 0 {
				fmt.Println("輸入不正確，請再選一次。")
				continue
			}
		} else {
			fmt.Println("沒有自動找到遊戲資料夾。")
		}
		p := strings.Trim(prompt(`請輸入 LIVE 資料夾路徑（例如 C:\Program Files\Roberts Space Industries\StarCitizen\LIVE），直接 Enter 取消：`), `" `)
		if p == "" {
			fmt.Println("已取消。")
			return ""
		}
		if isGameFolder(p) {
			return filepath.Clean(p)
		}
		if isGameFolder(filepath.Join(p, "LIVE")) {
			return filepath.Join(p, "LIVE")
		}
		fmt.Println("這個資料夾裡找不到 Data.p4k 或 Bin64\\StarCitizen.exe，請確認路徑。")
		list = nil
	}
}

// ---------- user.cfg ----------

var (
	reLang  = regexp.MustCompile(`(?i)^\s*g_language\s*=\s*(.*?)\s*$`)
	reAudio = regexp.MustCompile(`(?i)^\s*g_languageAudio\s*=`)
)

// 讀出 user.cfg 目前的 g_language（沒有就回傳空字串）
func currentLanguage(p string) string {
	b, err := os.ReadFile(p)
	if err != nil {
		return ""
	}
	cur := ""
	for _, l := range strings.Split(strings.ReplaceAll(strings.TrimPrefix(string(b), "\ufeff"), "\r\n", "\n"), "\n") {
		if m := reLang.FindStringSubmatch(l); m != nil {
			cur = m[1]
		}
	}
	return cur
}

// 設定 user.cfg 的 g_language，其他設定保留；addAudio 時若沒有語音設定就補上英文語音
func setLanguage(p, lang string, addAudio bool) error {
	var lines []string
	eol := "\r\n"
	if b, err := os.ReadFile(p); err == nil {
		t := strings.TrimPrefix(string(b), "\ufeff")
		if !strings.Contains(t, "\r\n") && strings.Contains(t, "\n") {
			eol = "\n"
		}
		t = strings.ReplaceAll(t, "\r\n", "\n")
		lines = strings.Split(strings.TrimRight(t, "\n"), "\n")
		if len(lines) == 1 && lines[0] == "" {
			lines = nil
		}
	}
	hasLang, hasAudio := false, false
	for i, l := range lines {
		if reLang.MatchString(l) {
			lines[i] = "g_language = " + lang
			hasLang = true
		} else if reAudio.MatchString(l) {
			hasAudio = true
		}
	}
	if !hasLang {
		lines = append(lines, "g_language = "+lang)
	}
	if addAudio && !hasAudio {
		lines = append(lines, "g_languageAudio = english")
	}
	return os.WriteFile(p, []byte(strings.Join(lines, eol)+eol), 0o644)
}

// ---------- 工具 ----------

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, b, 0o644)
}

func exeDir() string {
	if p, err := os.Executable(); err == nil {
		return filepath.Dir(p)
	}
	d, _ := os.Getwd()
	return d
}

func prompt(msg string) string {
	fmt.Print(msg)
	s, _ := in.ReadString('\n')
	return strings.TrimSpace(s)
}

func ask(msg string) bool {
	for {
		s := strings.ToUpper(prompt(msg))
		if s == "Y" || s == "YES" {
			return true
		}
		if s == "N" || s == "NO" || s == "" {
			return false
		}
	}
}

func fail(msg string) { fmt.Println("\n✘ " + msg) }

func pause() {
	if runtime.GOOS == "windows" {
		fmt.Println()
		prompt("按 Enter 關閉…")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
