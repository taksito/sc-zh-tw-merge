# 資料來源與致謝

本工具只負責**轉換與合併**，語言包的內容都來自以下作者與社群，著作權屬於各原作者。

## 語言包來源

### ExoAE — ScCompLangPack（Remix2）
- 作者：**ExoAE**
- 專案：<https://github.com/ExoAE/ScCompLangPack>
- 使用版本：[ScCompLangPackRemix2](https://github.com/ExoAE/ScCompLangPack/tree/main/ScCompLangPackRemix2)（`data/exocomp.ini`、`data/exocomp.base.ini`）
- 內容：零件類別／尺寸／等級前綴（例如 `MIL-2A "XL-1"`），以及下列 MrKraken 的改動
- Remix2 命名構想：[andydsm](https://www.reddit.com/user/andydsm/)；Remix 版本：[BeltaKoda](https://github.com/BeltaKoda/ScCompLangPackRemix)
- 原專案說明：「If you want to implement this in your own language pack, feel free to use the changes from the "modified_global.ini".」

### MrKraken — StarStrings
- 作者：**MrKraken**（[YouTube](https://www.youtube.com/@Mrkraken)・[X](https://x.com/RealMrKraken)・[Ko-fi](https://ko-fi.com/realmrkraken)）
- 專案：<https://github.com/MrKraken/StarStrings>
- 內容（經由 ExoAE 的語言包帶入）：
  - 任務標題的 `[BP]`／`[BP]*` 藍圖標記與聲望數值
  - 任務說明尾端的「Reputation Awarded」「Potential Blueprints」藍圖清單
  - 貨運任務標題的起點 → 終點
  - 劇情進度點數（Scenario Progress Points）
  - 礦物掃描的 RS 值、`Heph (Raw)` 等採礦文字縮短
  - 飛彈類型前綴（如 `[EM]`）、違禁品 `[!]` 標記
- ExoAE 的說明中表示 MrKraken 做了大部分的工作，如果想支持，請支持 MrKraken。

### CitizenWiki 中文社群 — 42kit 半漢化包
- 頁面：[CitizenWiki.cn 在地化（Localization）](https://citizenwiki.cn/index.php?title=Localization&variant=zh-tw)
- 檔案：<https://ini.42kit.com/half/global.ini>（2026-10-28 下載一次，存為 `data/chinese.ini`，之後不再自動抓取）
- 內容：簡體中文翻譯（本工具再轉為繁體）
- 如原作者不同意本專案保留此檔案，請開 Issue，會立即移除。

### 星際公民汉化组 — 飛船中文名稱
- 網站：<https://starcitizenzw.com/>
- 資料：<https://github.com/StarCitizenToolBox/LocalizationData>（`main` 分支「简中全汉化」4.10.2(PU)_CN_V2，2026-10-25）
- 使用範圍：只擷取 `vehicle_Name*` 的飛船名稱（628 條），存為 `src/ship_names_zh.txt`，由本工具轉成繁體並顯示成「中文 英文」
- 如原作者不同意此使用方式，請開 Issue，會立即移除。

### 中文輸入解碼表
- 檔案：`src/chinesedecode.txt`（7,020 條，`00=一` … `4EZ=龠`）
- 用途：讓遊戲顯示其他玩家輸入的中文
- 來源：由本專案維護者提供

## 第三方程式庫與資源

| 名稱 | 用途 | 授權 |
|---|---|---|
| [OpenCC](https://github.com/BYVoid/OpenCC) / [opencc-js](https://github.com/nk2028/opencc-js) | 簡體轉繁體 | Apache-2.0 |
| [JSZip](https://stuk.github.io/jszip/) | 產生 zip 下載檔 | MIT 或 GPLv3 |
| [IBM Plex Sans](https://github.com/IBM/plex)、[JetBrains Mono](https://github.com/JetBrains/JetBrainsMono)、[Noto Sans TC](https://fonts.google.com/noto/specimen/Noto+Sans+TC)（Google Fonts） | 網頁字型 | SIL Open Font License 1.1 |
| [Go](https://go.dev/) 標準函式庫 | 安裝器 | BSD-3-Clause |

## Cloud Imperium Games

- Star Citizen®、Roberts Space Industries® 與 Cloud Imperium® 是 Cloud Imperium Rights LLC 與 Cloud Imperium Rights Ltd. 的註冊商標。
- 遊戲文字的原始內容屬於 Cloud Imperium Games。
- 社群在地化政策：[Star Citizen: Community Localization Update](https://robertsspaceindustries.com/spectrum/community/SC/forum/1/thread/star-citizen-community-localization-update)（2023-10-11）
- [RSI Terms of Service](https://robertsspaceindustries.com/en/tos)
- [Translation & Fan Localization Statement](https://support.robertsspaceindustries.com/hc/en-us/articles/360006895793-Star-Citizen-Fankit-and-Fandom-FAQ#h_01JNKSPM7MRSB1WNBW6FGD2H98)
