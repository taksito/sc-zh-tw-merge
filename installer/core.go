// SC 繁中語言包安裝器：共用邏輯（找遊戲、讀語言包、備份、寫入、改 user.cfg）
package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	appName    = "SC 繁中語言包安裝器"
	version    = "2.0.0"
	langFolder = "chinese_(traditional)"
	backupDir  = "語言包備份"
)

// ---------- 語言包 ----------

type pack struct {
	Path    string
	Data    []byte
	Version string
	Lines   int
}

var rePackName = regexp.MustCompile(`(?i)^(global.*\.zip|global( \(\d+\))?\.ini(\.txt)?|global\.ini( \(\d+\))?\.txt)$`)

// 依序：命令列參數（拖曳到圖示上）→ 程式旁邊 → 下載資料夾；同一處取最新
func packCandidates(args []string) []string {
	out := append([]string{}, args...)
	dirs := []string{exeDir()}
	if home, err := os.UserHomeDir(); err == nil {
		dirs = append(dirs, filepath.Join(home, "Downloads"))
	}
	type fi struct {
		p string
		t time.Time
	}
	var found []fi
	for _, d := range dirs {
		ents, _ := os.ReadDir(d)
		for _, e := range ents {
			if e.IsDir() || !rePackName.MatchString(e.Name()) {
				continue
			}
			if info, err := e.Info(); err == nil {
				found = append(found, fi{filepath.Join(d, e.Name()), info.ModTime()})
			}
		}
	}
	sort.Slice(found, func(i, j int) bool { return found[i].t.After(found[j].t) })
	for _, f := range found {
		out = append(out, f.p)
	}
	return out
}

func findPack(args []string) (*pack, error) {
	for _, c := range packCandidates(args) {
		if p, err := loadPack(c); err == nil {
			return p, nil
		}
	}
	return nil, fmt.Errorf("找不到語言包。請從網站下載 global_zh-TW.zip，或按「瀏覽…」選擇檔案")
}

func loadPack(path string) (*pack, error) {
	var b []byte
	if strings.EqualFold(filepath.Ext(path), ".zip") {
		zr, err := zip.OpenReader(path)
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
				b, err = io.ReadAll(rc)
				rc.Close()
				if err != nil {
					return nil, err
				}
				break
			}
		}
		if b == nil {
			return nil, fmt.Errorf("zip 裡沒有 global.ini")
		}
	} else {
		var err error
		if b, err = os.ReadFile(path); err != nil {
			return nil, err
		}
	}
	if len(b) < 100000 || !bytes.Contains(b[:min(len(b), 4096)], []byte("=")) {
		return nil, fmt.Errorf("這個檔案不像是語言包")
	}
	v := "版本不明"
	if m := regexp.MustCompile(`(?m)^Frontend_PU_Version=(.*)$`).FindSubmatch(b); m != nil {
		v = strings.TrimSpace(string(m[1]))
	}
	return &pack{Path: path, Data: b, Version: v, Lines: bytes.Count(b, []byte("\n")) + 1}, nil
}

// ---------- 遊戲資料夾 ----------

func isGameFolder(p string) bool {
	for _, f := range []string{"Data.p4k", filepath.Join("Bin64", "StarCitizen.exe")} {
		if _, err := os.Stat(filepath.Join(p, f)); err == nil {
			return true
		}
	}
	return false
}

// 使用者選的可能是 StarCitizen 或 LIVE 本身
func normalizeGameFolder(p string) (string, bool) {
	p = filepath.Clean(strings.Trim(p, `" `))
	if isGameFolder(p) {
		return p, true
	}
	if isGameFolder(filepath.Join(p, "LIVE")) {
		return filepath.Join(p, "LIVE"), true
	}
	return p, false
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
	// RSI 啟動器紀錄裡的安裝路徑
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
	// 每個磁碟的常見安裝位置
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
	sort.SliceStable(out, func(i, j int) bool { return isLive(out[i]) && !isLive(out[j]) })
	return out
}

func isLive(p string) bool { return strings.EqualFold(filepath.Base(p), "LIVE") }

// ---------- user.cfg ----------

var (
	reLang  = regexp.MustCompile(`(?i)^\s*g_language\s*=\s*(.*?)\s*$`)
	reAudio = regexp.MustCompile(`(?i)^\s*g_languageAudio\s*=`)
)

func cfgPath(live string) string { return filepath.Join(live, "user.cfg") }
func targetPath(live string) string {
	return filepath.Join(live, "Data", "Localization", langFolder, "global.ini")
}

func currentLanguage(live string) string {
	b, err := os.ReadFile(cfgPath(live))
	if err != nil {
		return ""
	}
	cur := ""
	for _, l := range strings.Split(strings.ReplaceAll(strings.TrimPrefix(string(b), "\uFEFF"), "\r\n", "\n"), "\n") {
		if m := reLang.FindStringSubmatch(l); m != nil {
			cur = m[1]
		}
	}
	return cur
}

func isChineseInstalled(live string) bool {
	return strings.EqualFold(currentLanguage(live), langFolder)
}

func setLanguage(live, lang string, addAudio bool) error {
	p := cfgPath(live)
	var lines []string
	eol := "\r\n"
	if b, err := os.ReadFile(p); err == nil {
		t := strings.TrimPrefix(string(b), "\uFEFF")
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

// ---------- 備份、安裝、移除 ----------

func backupRoot(live string) string { return filepath.Join(live, backupDir) }

func backup(live string, files ...string) (string, int, error) {
	stamp := time.Now().Format("20060102-150405")
	bdir := filepath.Join(backupRoot(live), stamp)
	for i := 2; ; i++ { // 同一秒執行兩次時不要蓋掉前一份
		if _, err := os.Stat(bdir); err != nil {
			break
		}
		bdir = filepath.Join(backupRoot(live), stamp+"-"+strconv.Itoa(i))
	}
	n := 0
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			continue
		}
		rel, _ := filepath.Rel(live, f)
		if err := copyFile(f, filepath.Join(bdir, rel)); err != nil {
			return bdir, n, err
		}
		n++
	}
	return bdir, n, nil
}

const adminHint = "（遊戲裝在 Program Files 時，可能要對安裝器按右鍵「以系統管理員身分執行」）"

// 安裝或更新中文化；log 每一步回報進度
func install(live string, p *pack, log func(string)) error {
	target, cfg := targetPath(live), cfgPath(live)
	bdir, n, err := backup(live, target, cfg)
	if err != nil {
		return fmt.Errorf("備份失敗：%v\n%s", err, adminHint)
	}
	if n > 0 {
		log(fmt.Sprintf("✔ 已備份 %d 個舊檔案到 %s", n, bdir))
	} else {
		log("✔ 沒有舊檔案需要備份")
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("建立資料夾失敗：%v\n%s", err, adminHint)
	}
	data := p.Data
	if !bytes.HasPrefix(data, []byte{0xEF, 0xBB, 0xBF}) {
		data = append([]byte{0xEF, 0xBB, 0xBF}, data...)
	}
	if err := os.WriteFile(target, data, 0o644); err != nil {
		return fmt.Errorf("寫入 global.ini 失敗：%v\n%s", err, adminHint)
	}
	log("✔ 已寫入 " + target)
	if err := setLanguage(live, langFolder, true); err != nil {
		return fmt.Errorf("修改 user.cfg 失敗：%v\n%s", err, adminHint)
	}
	log("✔ 已把 user.cfg 設為 g_language = " + langFolder)
	log("完成！啟動遊戲就會是繁體中文。")
	return nil
}

// 移除中文化：只把 user.cfg 改回英文，中文 global.ini 保留
func uninstall(live string, log func(string)) error {
	bdir, n, err := backup(live, cfgPath(live))
	if err != nil {
		return fmt.Errorf("備份失敗：%v\n%s", err, adminHint)
	}
	if n > 0 {
		log(fmt.Sprintf("✔ 已備份 user.cfg 到 %s", bdir))
	}
	if err := setLanguage(live, "english", false); err != nil {
		return fmt.Errorf("修改 user.cfg 失敗：%v\n%s", err, adminHint)
	}
	log("✔ 已把 user.cfg 設為 g_language = english")
	log("完成！啟動遊戲就會是英文。中文語言包仍保留，之後可隨時再安裝回來。")
	return nil
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
