# SC 繁中語言包合併台

把 Star Citizen 的**簡體中文漢化包**與 **ExoAE 英文語言包**合併成**繁體中文 `global.ini`**（UTF-8 with BOM）。

整個工具是一個網頁，在瀏覽器裡就能跑，不用安裝任何東西；檔案不會上傳到任何伺服器。

## 功能

- **簡轉繁**：使用 [OpenCC](https://github.com/nk2028/opencc-js)，可選台灣正體＋台灣用語、只轉字、香港繁體。
- **接上 ExoAE 的資訊**
  - 零件等級：`MIL-1C "Bracer"` → `軍規1C-Bracer`（MIL 軍規、CIV 民用、IND 工業、STL 隱形、CMP 競速）
  - 任務標題：`[100 聲望] [藍圖]*`、貨運路線
  - 任務說明：獲得聲望、可能獲得的藍圖清單
- **礦石壓縮**：「中文＋英文縮寫（3–6 字母）＋RS 值」，每行不超過設定字數（預設 14）。
- **自動檢查**：自動合併後會檢查變數是否缺漏、`<EM4>` 是否成對；有問題的另外列出。
- **審核介面**：只需看「問題項目」「需人工確認」「礦石」，其餘自動套用。
- **進度自動儲存**：存在瀏覽器（localStorage＋IndexedDB），也可匯出／匯入審核紀錄 JSON。
- **中文輸入解碼表**：輸出的 `global.ini` 最後自動加上 7,020 條對照（`00=一` … `4EZ=龠`，來源 `src/chinesedecode.txt`），讓遊戲能顯示其他玩家輸入的中文。可在匯出區取消勾選。
- **選填原版英文**：放入同版本原版 `global.ini`，可精確找出 ExoAE 改了哪些條目。

## 使用

網址：https://taksito.github.io/sc-zh-tw-merge/

1. 按「從網路取得最新版並合併」，會讀網站上的：
   - 漢化：`data/chinese.ini`，固定版本（2026-10-28 從 [42kit 半漢化包](https://citizenwiki.cn/index.php?title=Localization&variant=zh-tw) 下載一次，之後不再自動抓）
   - ExoAE：`data/exocomp.ini`，[ScCompLangPackRemix2](https://github.com/ExoAE/ScCompLangPack/tree/main/ScCompLangPackRemix2) 的最新版
   
   也可以手動選擇檔案再按「開始合併」。
2. ExoAE 改版後，工具會拿 `data/exocomp.base.ini`（漢化當時對照的版本）比對：
   - **英文有更新**：英文變了，中文可能過時，請確認
   - **新增條目**：新版才有、漢化沒有的條目
3. 看完左側「需要你看的」分類。
4. 按「下載 global.ini」。
5. 放到 `StarCitizen\LIVE\Data\Localization\chinese_(traditional)\global.ini`，並在 `LIVE\user.cfg` 加上：
   ```
   g_language = chinese_(traditional)
   ```

## 一鍵安裝到遊戲（Windows）

下載 [`download/SC-zhTW-Installer.exe`](https://taksito.github.io/sc-zh-tw-merge/download/SC-zhTW-Installer.exe) 後執行：

1. 自動找語言包：拖曳到程式上的檔案 → 程式旁邊 → 「下載」資料夾裡最新的 `global*.zip` / `global.ini`
2. 自動找遊戲：讀 RSI 啟動器紀錄＋各磁碟常見安裝位置，LIVE 優先；找不到可手動輸入
3. 確認後寫入 `LIVE\Data\Localization\chinese_(traditional)\global.ini`（UTF-8 with BOM），並把 `user.cfg` 的 `g_language` 設為 `chinese_(traditional)`（其他設定保留）
4. 舊的 `global.ini`、`user.cfg` 先備份到 `LIVE\語言包備份\日期時間\`

已經安裝過中文時（`user.cfg` 的 `g_language` 是 `chinese_(traditional)`），會先問要做哪一件事：

- **更新中文化**：照上面步驟寫入新的語言包
- **移除中文化（變回英文）**：備份 `user.cfg` 後只把 `g_language` 改成 `english`，中文 `global.ini` 保留，之後可隨時再換回中文

原始碼在 `installer/`（Go，只用標準函式庫）；修改後 `.github/workflows/build-installer.yml` 會自動重新編譯。

## 自動更新

`.github/workflows/update-sources.yml` 每天台灣時間 04:17 檢查 ExoAE Remix2，有新版才更新 `data/exocomp.ini` 與 `data/sources.json`。
漢化不會自動更新。要立刻檢查：Actions →「檢查 ExoAE 更新」→「Run workflow」。

## 共用修改（審核紀錄）

在工具裡按「匯出審核紀錄」得到 `sc-merge-review.json`，改名為 `review.json` 上傳到 `data/`。
之後任何電腦打開網站都會自動套用這些修改（本機已有的修改優先）。

## 開發

```
src/engine.js   合併邏輯（瀏覽器與 Node 共用）
src/app.html    介面
build.py        組成單一檔案 index.html
test/merge-test.js  用 Node 檢查合併結果
```

修改 `src/` 之後執行：

```
python3 build.py
node test/merge-test.js chinese.ini exocomp.ini
```

## 啟用 GitHub Pages

Repository → Settings → Pages → Source 選 **Deploy from a branch**，Branch 選 `main`、資料夾 `/ (root)`。

## 資料來源與致謝

- **ExoAE** — [ScCompLangPack](https://github.com/ExoAE/ScCompLangPack)（使用 Remix2）：零件類別／尺寸／等級前綴
- **MrKraken** — [StarStrings](https://github.com/MrKraken/StarStrings)：任務藍圖與聲望、貨運路線、採礦 RS 值等改動（經由 ExoAE 語言包帶入）
- **CitizenWiki 中文社群（42kit 半漢化包）** — [CitizenWiki.cn 在地化](https://citizenwiki.cn/index.php?title=Localization&variant=zh-tw)：簡體中文翻譯
- **OpenCC / opencc-js**（Apache-2.0）：簡轉繁；**JSZip**（MIT）：產生 zip

完整清單見 **[CREDITS.md](CREDITS.md)**。

## 使用者協議與授權

- 使用前請閱讀 **[使用者協議（TERMS.md）](TERMS.md)**。
- 本專案自己撰寫的程式碼採用 [MIT 授權](LICENSE)；`data/` 內的語言包與 `src/chinesedecode.txt` 不在此授權內，權利屬於各原作者。
- 這是非官方的 Star Citizen 粉絲專案，與 Cloud Imperium Games 無關。Star Citizen® 為 Cloud Imperium Rights LLC 的註冊商標。
- CIG 的社群在地化政策：[Star Citizen: Community Localization Update](https://robertsspaceindustries.com/spectrum/community/SC/forum/1/thread/star-citizen-community-localization-update)
