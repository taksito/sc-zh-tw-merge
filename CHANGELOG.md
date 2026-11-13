# 更新紀錄

## 2026-11-13
- **4.10.3 正式版**：ExoAE、MrKraken StarStrings 更新到 4.10.3。正式版新增 281 條、英文有更新 41 條全部補上繁中，「需要你看」回到 0。補充翻譯檔改名為 `src/extra_zh.txt`，正式版與 PTU 共用（只套用在來源裡有的鍵）。物品名稱照慣例保留英文。
- **預計出現的艦船**：任務說明的艦船清單標籤（Targets、Hostiles、Reinforcements、Support Ships、Scouts、Defenders、Bombers 等）與通用艦種（例如 Hard Criminal Large Combat Ship → 不法分子大型戰鬥艦（困難））自動翻譯；StarStrings 的 Cooldown、Debug Name 視為附加資訊並翻譯。
- 「英文有更新」改為忽略 `<EM>` 強調標記與空白差異。

## 2026-11-01
- **聲望等級加上 LV**：`RepScope_*_Rank*`、`RepStanding_*_Rank*`、`RacingRep_Ship_Rank*` 的等級名稱前加上 LV.1～LV.N（例：LV.1 初級承包商），第 0 級（申請者／中立）不加。「自訂介面」新增「聲望等級」分類，設定裡可取消勾選。
- **PTU 頁面調整**：PTU 版的說明、分類與卡片不再寫 ExoAE，改標示英文來源為 MrKraken StarStrings。「自訂介面」新增「PTU 新增（已翻譯）」「PTU 更新（已翻譯）」兩個分類，可逐條查看英文原文與翻譯；更新的條目會同時顯示正式版英文。
- **PTU 條目翻譯**：PTU 新增 95 條、英文有更新 24 條全部翻成繁中（`src/ptu_zh.txt`，只用於 PTU 版）。包含尼克斯星系互安防務任務、登船救援、斯卡諾夫步槍、奧馬洛夫瞄準鏡、劫掠者迷彩裝備等；ShipSelector 的「飛船」改為「載具」。物品名稱照正式版慣例保留英文。PTU 移除的 12 條仍保留在輸出裡。
- 任務說明的「Ships to Expect」波次清單視為附加資訊，自動翻成「預計出現的艦船／第 N 波」，通用艦種（例如 Medium Criminal Fighter）也會翻譯。
- **新增 PTU 搶先體驗版**：頁面上方可切換「LIVE 正式版／PTU 搶先體驗版」。PTU 的英文來自 MrKraken StarStrings `sc-ptu`（4.10.3 PTU），對照正式版找出新增 95 條、英文有更新 24 條，先顯示英文待翻譯。
- 「英文有更新」只看原文本體，藍圖清單或聲望標籤的變動不再列入。

## 2026-10-31
- **萬向節顯示修正**：HUD 武器瞄準區 `hud_aimable_*` 在遊戲裡會變成亂碼，全部改用英文：萬向節狀態 FXD、MAN（Fixed、Unlocked），瞄準方式 AUTO、PIP、PRE（Auto、Default、Precision）。
- 萬向節模式縮寫 `hud_gimbal_mode_short_*` 也維持英文（AUTO、FXD、MAN、FXD-PRE、MAN-PRE）。
- 「固定用 ExoAE 原文」清單支援結尾 `*`，代表該開頭的所有條目。
- 修正：之前開過網頁的人，瀏覽器存的舊詞彙表會擋住新版預設（例如這次的萬向節、AAT-34 砲台）。現在會自動補上新的預設條目，你自己的修改保留。

## 2026-10-29
- 任務地點（`mission_location*`）加上中文：英文 3 個詞以上只用中文，1～2 個詞顯示「英文 中文」，PYAM 等縮寫保留。
- 飛船名稱改為中文一行、英文下一行。

## 2026-10-28
- 飛船名稱加上中文（中文在前）。
- AAT-34、AAT-21 Turret 改為砲台；補上漏翻的任務標記（量子信標等）。
- 地點名稱改為「英文 中文」。
- 新增手動安裝說明與「直接下載 global.ini」。
- 安裝器改為視窗版（v2.0.0），可更新或移除中文化。
- 輸出的 global.ini 附加中文輸入解碼表。
- 首頁改為三步驟懶人包。
