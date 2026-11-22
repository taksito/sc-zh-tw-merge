# 更新紀錄

## 2026-11-22
- **全繁中版：星圖地名全部中文化**：合併台用的漢化組「半漢化」包會故意把地面名稱留英文，所以 F2 星圖上像 Seraphim Station、Cellin、Daymar、CRU L1 這些都還是英文。現在從漢化組「全漢化」4.10.3(PU)_CN_V0.1（2026-11-15）補上 1,922 條：熾天使空間站、熾天使診所、賽琳、戴瑪爾、奧里森、十 CRU L1、殁氧星、端點星、各前哨站與地下設施、拉格朗日休息站，以及勢力名稱、人名、任務物品等。只在全繁中版的輸出仍是英文時才套用，原本是中文的條目和雙語版都不受影響。
  - 新增 `src/full_names_zh.txt`（由 `tools/make_full_names.py` 產生），變數記號不一致的條目自動跳過。
  - 修正：檔案裡帶 `,P` 後綴的鍵（約 1.3 萬條）之前查不到全繁中版的對照表。
  - 仍是英文的只剩漢化組本身也沒翻的代號（NT-999-XV、小行星帶編號 2904-UEE649 等）。

## 2026-11-21
- **補上 Idris 雷達名稱**：遊戲裡會顯示 `@item_Name_Radr_GNRP_S03_idris`，因為官方英文檔本身就漏了這條（鍵名還把 GRNP 拼成 GNRP），ExoAE、StarStrings、汉化组也都沒有。現在補成「伊德里斯 雷達」（英文 Idris Radar，比照 Mauler Radar）。
- `live_missing_*.txt` 新增「手動補充」區，`tools/make_live_missing.py` 重新產生時會保留。

## 2026-11-19
- **全繁中版跟上雙語版**：逐條比對兩版（98,396 條，鍵完全一致），補上全繁中版落後的地方：
  - 漢化包本身留英文的 19 個地名（Hurston、ArcCorp、Area18、Grim HEX、Levski、星門、出生地點選單等），雙語版有「英文 中文」、全繁中版卻只有英文，改用地點對照表的中文（赫斯頓、18區、六角灣…）。
  - 藍圖清單的燃料噴嘴（Marlin 槍魚、Lindstrom 林斯特龍、Bendix 奔德士、Torrez 托雷斯、Ezra 以斯拉）鍵名是 `Nozzle_…_Name`，之前沒查到，改為一起查。
  - 預計出現的艦船清單裡的系列名稱（Anvil Hornet（任意型號）→ 鐵砧 大黃蜂、RSI Aurora Mk I → 極光、C.O. Mustang、Constellation Mk IV、C8 Pisces）改用各型號共同的中文。
  - 通用艦種補上 Heavy Fighter、Medium Fighter、Heavy Interceptor、Distortion Ship、Vanduul Combat Ship（雙語版也一起翻）。
  - 仍是英文的只剩 Anvil Paladin（漢化組沒有中文名）。
- **飛船改成英文在前**：船名改為「英文 ⏎ 中文」（例：`Drake Cutlass Black` ⏎ `德雷克 黑彎刀`），遊戲裡的飛船清單照廠商英文名排序（Aegis → Anvil → Argo → … → Vanduul）。完整船名缺廠商的自動補上（`Cutlass Black PYAM Exec` → `Drake Cutlass Black PYAM Exec`、`Aurora Mk II TS Module` → `RSI Aurora Mk II TS Module`、`ARGO MOLE Alliance` → `Argo MOLE Alliance`）；簡稱（`_short`）不加廠商。設定新增「飛船順序」，可改回中文在前。全繁中版不受影響。

## 2026-11-18
- **補漏條目加上中文**：146 條補漏中仍只有英文的 6 條（`event_ShipName_ATLS_GEO`、`event_ShipName_Hyun`、`event_ShipName_UTV`、`event_ShipTitle_S65`、`vehicle_NameDRAK_Marauder_short,P`、`vehicle_NameGAMA_Hyun_Short`）在 `src/extra_zh.txt` 加上中文。
- **補上 ExoAE 漏掉的 146 條**：ExoAE 4.10.3（pre release）缺少新飛船 Drake Marauder（`vehicle_NameDRAK_Marauder` 等）與 IAE 2956 活動條目，遊戲會顯示鍵名。正式版改用 MrKraken StarStrings 正式版的英文補上（`src/live_missing_en.txt`），中文取自星際公民汉化组 4.10.3 LIVE 汉化临时版 V0.1（`src/live_missing_zh_s.txt`），由 `tools/make_live_missing.py` 產生。ExoAE 補上後會自動改用 ExoAE。船名加入 `src/ship_names_zh.txt`。

## 2026-11-16
- **全繁中版加回聲望與獎勵**：任務標題的聲望／藍圖標籤與說明尾端的獲得聲望、可能獲得的藍圖、預計出現的艦船都保留；藍圖清單的物品／零件名稱與艦船清單的船名換成中文（對照汉化组全汉化名稱），括號內的零件類型（Shield、Cooler、Quantum Drive…）與彈匣容量（N cap）也翻成中文。

## 2026-11-14
- **修正 `@item_Nameksar_lmg_energy_01`**：4.10.3 新武器 Verdict LMG（裁決輕機槍）基本款的名稱、簡稱、說明在 ExoAE、StarStrings、漢化包都缺，遊戲會顯示鍵名。`src/extra_zh.txt` 新增「+鍵名」寫法，可強制補上所有語言包都缺的鍵；全繁中版名稱為「裁決 輕機槍」。

## 2026-11-13
- **全繁中版：物品也翻中文**：物品、武器、裝備、零件與塗裝名稱改用星際公民汉化组「简中全汉化」的中文名稱（`src/item_names_full_zh.txt`，約 1 萬 600 條，由 `tools/make_item_names.py` 產生）。4.10.2 全汉化為主，4.10.3 英文有改的鍵用 5.0 PTU 全汉化。零件保留「軍規1B-」等級標示。只影響全繁中版。
- **全繁中版**：網頁三個步驟下方新增「全繁中版」下載（`global_zh-TW_pure.zip`／直接下載 global.ini）。只用中文，不加英文船名、英文地名、礦石縮寫、任務標題標籤與 ExoAE 附加資訊；新條目用補充翻譯。HUD 萬向節字串為避免亂碼維持英文。
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
