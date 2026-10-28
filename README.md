# SC 繁中語言包合併台

把 Star Citizen 的**簡體中文漢化包**與 **ExoComp 英文語言包**合併成**繁體中文 `global.ini`**（UTF-8 with BOM）。

整個工具是一個網頁，在瀏覽器裡就能跑，不用安裝任何東西；檔案不會上傳到任何伺服器。

## 功能

- **簡轉繁**：使用 [OpenCC](https://github.com/nk2028/opencc-js)，可選台灣正體＋台灣用語、只轉字、香港繁體。
- **接上 ExoComp 的資訊**
  - 零件等級：`MIL-1C "Bracer"` → `軍規1C-Bracer`（MIL 軍規、CIV 民用、IND 工業、STL 隱形、CMP 競速）
  - 任務標題：`[100 聲望] [藍圖]*`、貨運路線
  - 任務說明：獲得聲望、可能獲得的藍圖清單
- **礦石壓縮**：「中文＋英文縮寫（3–6 字母）＋RS 值」，每行不超過設定字數（預設 14）。
- **自動檢查**：自動合併後會檢查變數是否缺漏、`<EM4>` 是否成對；有問題的另外列出。
- **審核介面**：只需看「問題項目」「需人工確認」「礦石」，其餘自動套用。
- **進度自動儲存**：存在瀏覽器（localStorage＋IndexedDB），也可匯出／匯入審核紀錄 JSON。
- **選填原版英文**：放入同版本原版 `global.ini`，可精確找出 ExoComp 改了哪些條目。

## 使用

1. 打開 `index.html`（或 GitHub Pages 網址）。
2. 選擇簡中漢化包（如 `chinese.ini`）和 ExoComp 英文包，按「開始合併」。
3. 看完左側「需要你看的」分類。
4. 按「下載 global.ini」。
5. 放到 `StarCitizen\LIVE\Data\Localization\chinese_(traditional)\global.ini`，並在 `LIVE\user.cfg` 加上：
   ```
   g_language = chinese_(traditional)
   ```

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

## 授權

工具本身程式碼可自由使用。語言包內容的著作權屬於原作者（Cloud Imperium Games、漢化團隊、ExoComp），本專案不包含任何語言包檔案。
