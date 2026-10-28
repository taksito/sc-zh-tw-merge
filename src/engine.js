// ===== SC 語言包合併引擎（瀏覽器與 Node 共用）=====
(function (root) {
  const CJK = /[\u3400-\u9fff\uf900-\ufaff]/;
  const KEY_LINE = /^[^\s=\\][^=]*?=/; // 鍵不以反斜線開頭；值內的斷行殘段（如 "\nCarrying..."）會接回上一行

  function parseIni(text) {
    if (text.charCodeAt(0) === 0xfeff) text = text.slice(1);
    const crlf = text.indexOf('\r\n') !== -1;
    text = text.replace(/\r\n?/g, '\n');
    const keys = [], map = new Map(), repaired = new Set(), dups = [];
    let last = null;
    for (const line of text.split('\n')) {
      if (line === '') continue;
      const m = KEY_LINE.exec(line);
      if (m && !/\s/.test(m[0].slice(0, -1))) {
        const k = m[0].slice(0, -1), v = line.slice(m[0].length);
        if (map.has(k)) dups.push(k); else keys.push(k);
        map.set(k, v); last = k;
      } else if (last !== null) {
        map.set(last, map.get(last) + line); repaired.add(last);
      }
    }
    return { keys, map, crlf, repaired: [...repaired], dups };
  }

  // ---------- 字典 ----------
  // 每行： 英文 => 中文 ；以 / 開頭結尾的是正規式（可用 $1）；# 開頭為註解
  const DEFAULT_TAIL_DICT = `# ExoAE 附加資訊的翻譯（只套用在從 ExoAE 接過來的片段）
Reputation Awarded (by difficulty): => 獲得聲望（依難度）：
Reputation Awarded: => 獲得聲望：
Multiple Blueprint Pools => 多個藍圖池
/Potential Blueprints \\((.+?) [Oo]nly\\)/ => 可能獲得的藍圖（僅限 {$1}）
Potential Blueprints => 可能獲得的藍圖
/Awarded from (.+?) level variants/ => {$1} 等級的任務變體可獲得
/\\bPool (\\d+)/ => 藍圖池 $1
Scenario Progress Points => 劇情進度點數
Total => 總計
Transport => 運輸
Defense => 防禦
Collection => 收集
(Split) => （分攤）
/\\[Regional Variants\\] example locations: ([^<]+)/ => [區域變體] 範例地點：{,$1}
/ Rep\\]/ =>  聲望]
[BP] => [藍圖]
Warning: Some blueprint strings in this pool may have been assigned to the wrong item by CIG. Check scmdb.net for manually corrected pools => 警告：此藍圖池中部分藍圖可能被 CIG 對應到錯誤物品，修正後的清單請至 scmdb.net 查詢`;

  const DEFAULT_TERM_DICT = `# 簡轉繁之後，再套用到所有中文內容的自訂替換（左邊寫繁體轉換後的結果）
# 例：
# 視頻 => 影片
`;

  function parseDict(txt) {
    const out = [];
    for (const raw of (txt || '').split('\n')) {
      const line = raw.replace(/\s+$/, '');
      if (!line.trim() || line.trim().startsWith('#')) continue;
      const i = line.indexOf(' => ');
      if (i < 0) continue;
      const a = line.slice(0, i), b = line.slice(i + 4);
      if (a.length > 2 && a.startsWith('/') && a.endsWith('/')) {
        try { out.push({ re: new RegExp(a.slice(1, -1), 'g'), to: b }); } catch (e) { /* 忽略錯誤正規式 */ }
      } else out.push({ s: a, to: b });
    }
    // 純文字依長度由長到短，避免短詞先吃掉長詞；正規式保持書寫順序並排在前面
    const res = out.filter(x => x.re), str = out.filter(x => !x.re).sort((x, y) => y.s.length - x.s.length);
    return res.concat(str);
  }

  function applyDict(s, dict, gloss) {
    for (const d of dict) {
      if (d.re) {
        s = s.replace(d.re, (...m) => {
          let t = d.to;
          t = t.replace(/\{,\$(\d)\}/g, (_, n) => (m[+n] || '').split(/,\s*/).map(x => glossLookup(x.trim(), gloss)).join('、'));
          t = t.replace(/\{\$(\d)\}/g, (_, n) => glossLookup(m[+n] || '', gloss));
          t = t.replace(/\$(\d)/g, (_, n) => m[+n] || '');
          return t;
        });
      } else if (s.includes(d.s)) s = s.split(d.s).join(d.to);
    }
    return s;
  }
  function glossLookup(en, gloss) { return (gloss && gloss.get(en)) || en; }

  // 從兩個檔案自動建立「短英文 → 中文」詞彙表（例如 Jr. Contractor → 初級承包商）
  function buildGlossary(zhT, exo) {
    const votes = new Map();
    for (const [k, ev] of exo.map) {
      if (ev.length > 40 || CJK.test(ev) || /[<~\\]/.test(ev)) continue;
      const zv = zhT.get(k);
      if (!zv || !CJK.test(zv) || zv.length > 30) continue;
      let m = votes.get(ev); if (!m) votes.set(ev, m = new Map());
      m.set(zv, (m.get(zv) || 0) + 1);
    }
    const g = new Map();
    for (const [en, m] of votes) g.set(en, [...m].sort((a, b) => b[1] - a[1])[0][0]);
    return g;
  }

  // ---------- 規則 ----------
  const P = '\\n\\n'; // ini 內的段落分隔（字面上的 \n\n）
  const normPunct = s => s.replace(/：/g, ':').replace(/，/g, ',').replace(/（/g, '(').replace(/）/g, ')').replace(/\s+/g, '');
  const trimTail = s => s.replace(/(?:\\n|\s)+$/, '');
  const TITLE_SUFFIX = /\s*(<EM4>(?:\s*\[[^\]<]*\]\*?)+\s*<\/EM4>)\s*$/;

  const EXO_MARK = /(?:\\n)+(?=<EM4>(?:Reputation Awarded|Potential Blueprints|Multiple Blueprint Pools|Scenario Progress Points|Pool \d|Awarded from))/;
  function emTailCount(v) {
    const ps = trimTail(v).split(P); let n = 0;
    for (let i = ps.length - 1; i >= 0; i--) { if (/^(?:\s|\\n)*<EM4>/.test(ps[i])) n++; else break; }
    return { n, ps };
  }
  const BR = /\[[^\]\n]{1,40}\]/g;
  function tagDiff(ev, zv) {
    // 只比「數量」：原版就有的 [Recon] 被翻成 [侦察] 不算差異；[ Groaning ] 這類語音標記忽略
    const cnt = (v, re) => (v.match(re) || []).filter(t => !/^\[ /.test(t)).length;
    const d = [];
    const be = cnt(ev, BR), bz = cnt(zv, BR) + cnt(zv, /【[^】]{1,40}】/g);
    if (be > bz) d.push('[標籤] ×' + (be - bz));
    const ce = cnt(ev, /<EM\d>/g), cz = cnt(zv, /<EM\d>/g);
    if (ce > cz) d.push('<EM> ×' + (ce - cz));
    return d;
  }

  const ROUTE = /\s*\|\s*((?:from |to )?~mission\([^)]*\)(?:\s*>\s*~mission\([^)]*\))?)\s*$/;
  const AUTO_CATS = new Set(['title', 'tail', 'wrap']);
  const toks = v => (v || '').match(/~\w+\([^)]*\)|%l?[sdif]/g) || [];
  function validate(e) {
    const iss = [], out = e.out || '';
    const open = (out.match(/<EM\d>/g) || []).length, close = (out.match(/<\/EM\d>/g) || []).length;
    const bal = v => (v.match(/<EM\d>/g) || []).length - (v.match(/<\/EM\d>/g) || []).length;
    if (open !== close && bal(out) !== bal(e.zh || '')) iss.push('<EM4> 標籤沒有成對');
    const have = new Set(toks(out));
    const miss = [...new Set(toks(e.exo))].filter(t => !have.has(t));
    if (miss.length) iss.push('少了變數 ' + miss.slice(0, 3).join(' '));
    // 從 ExoAE 接進來、翻譯後仍殘留的英文（物品清單「- xxx」和地名除外）
    const tail = out.slice(Math.max(0, out.length - (e.added || '').length - 20));
    const left = tail.split('\\n').filter(l => !/^\s*- /.test(l))
      .map(l => l.replace(/<\/?EM\d>/g, '').replace(/~\w+\([^)]*\)/g, '').replace(/\[[^\]]*\]/g, ''))
      .filter(l => /[A-Za-z]{4,}(?:\s+[A-Za-z]{2,}){1,}/.test(l) && !/[\u4e00-\u9fff]/.test(l.replace(/[A-Za-z].*$/, '')) );
    if (left.length) iss.push('附加資訊有英文沒翻到：' + left[0].trim().slice(0, 40));
    if (!e.zh || !e.zh.trim()) iss.push('繁中原譯是空的');
    return iss;
  }

  /**
   * 回傳 { entries: [...], stats }
   * entry: { k, cat, out, zh, zhRaw, exo, van, added }
   * cat：
   *   same      兩邊一樣或只需繁化（不審核）
   *   zh        採用繁中（不審核）
   *   en_exo    繁中檔仍是英文 → 採用 ExoAE（含標籤）
   *   title     任務標題補上 ExoAE 的 [聲望/藍圖] 標籤
   *   tail      任務說明補上 ExoAE 的獎勵/藍圖清單
   *   wrap      （有原版時）ExoAE 在原文前後加的內容，接到繁中
   *   check     ExoAE 有改動但無法自動接 → 需人工
   *   zh_only / exo_only  只存在於一邊
   *   updated   （有 opts.exoBase 時）ExoAE 英文跟漢化當時對照的版本不同 → 中文可能過時
   */
  function merge(zhRaw, zhT, exo, van, opts) {
    opts = opts || {};
    const tailDict = parseDict(opts.tailDict != null ? opts.tailDict : DEFAULT_TAIL_DICT);
    const gloss = buildGlossary(zhT, exo);
    const tr = s => applyDict(s, tailDict, gloss);
    const classMap = Object.assign({}, DEFAULT_CLASS_MAP, opts.classMap || {});
    const fill = new Map();
    for (const d of parseDict(opts.fill != null ? opts.fill : DEFAULT_FILL)) if (d.s) fill.set(d.s.trim(), d.to);
    const forceExo = new Set((opts.forceExo != null ? opts.forceExo : DEFAULT_FORCE_EXO).split(/\r?\n/).map(x => x.replace(/#.*/, '').trim()).filter(Boolean));
    const entries = [], stats = {};
    const add = (e) => { entries.push(e); stats[e.cat] = (stats[e.cat] || 0) + 1; };

    for (const k of exo.keys) {
      const ev = exo.map.get(k), zr = zhRaw.map.get(k), zv = zhT.get(k);
      const vv = van ? van.map.get(k) : undefined;
      const base = { k, zh: zv, zhRaw: zr, exo: ev, van: vv };
      if (forceExo.has(k)) { add({ ...base, cat: 'force', out: ev }); continue; }
      if (fill.has(k)) { add({ ...base, cat: 'fill', out: fill.get(k) }); continue; }
      if (zr === undefined) { add({ ...base, cat: 'exo_only', out: ev }); continue; }
      if (zv === ev) { add({ ...base, cat: 'same', out: zv }); continue; }

      const zhIsCJK = CJK.test(zr);
      const exoChanged = vv === undefined ? null : ev !== vv;

      if (!zhIsCJK && (normPunct(zr) === normPunct(ev) || exoChanged === false)) { add({ ...base, cat: 'same', out: zv }); continue; }
      if (zhIsCJK && exoChanged === false) { add({ ...base, cat: 'zh', out: zv }); continue; }

      // 有原版英文：精確找出 ExoAE 在原文前/後加了什麼
      if (vv !== undefined && vv.length > 0) {
        const at = ev.indexOf(vv);
        if (at >= 0) {
          const pre = ev.slice(0, at), post = ev.slice(at + vv.length);
          const out = (pre ? tr(pre) : '') + (post ? trimTail(zv) + tr(post) : zv);
          add({ ...base, cat: 'wrap', out, added: pre + (pre && post ? ' … ' : '') + post });
          continue;
        }
      }

      // 規則 1：任務標題結尾的 <EM4>[100 Rep] [BP]*</EM4>
      const mt = TITLE_SUFFIX.exec(ev);
      if (mt && !/<\/EM4>\s*$/.test(zv)) {
        // 貨運標題中間還有「| 起點 > 終點」路線，漢化沒有的話一起接上
        const body = ev.slice(0, mt.index);
        const mr = ROUTE.exec(body);
        let route = '';
        if (mr && !toks(mr[1]).every(t => zv.includes(t))) route = ' | ' + mr[1].replace(/^from\s+/, '自 ').replace(/^to\s+/, '至 ');
        add({ ...base, cat: 'title', out: zv.replace(/\s+$/, '') + route + ' ' + tr(mt[1]), added: (mr && route ? mr[0].trim() + ' ' : '') + mt[1] });
        continue;
      }
      // 規則 2a：ExoAE 的獎勵區塊（從第一個「Reputation Awarded / Potential Blueprints…」段落到結尾）
      const mk = EXO_MARK.exec(ev);
      if (mk && !EXO_MARK.test(zr)) {
        const extra = ev.slice(mk.index).replace(/^(?:\s|\\n)+/, '').replace(/(?:\\n|\s)+$/, '');
        add({ ...base, cat: 'tail', out: trimTail(zv) + P + tr(extra), added: extra });
        continue;
      }
      // 規則 2b：說明結尾多出的其他 <EM4> 段落
      const te = emTailCount(ev), tz = emTailCount(zv);
      if (te.n > tz.n && te.n - tz.n < te.ps.length) {
        const extra = te.ps.slice(te.ps.length - (te.n - tz.n)).join(P).replace(/^(?:\s|\\n)+/, '');
        add({ ...base, cat: 'tail', out: trimTail(zv) + P + tr(extra), added: extra });
        continue;
      }
      // 漢化本來就是英文的條目：零件等級翻成中文，其餘直接用 ExoAE
      if (!zhIsCJK) {
        const mc = COMP_RE.exec(ev);
        if (mc && classMap[mc[1].toUpperCase()]) {
          add({ ...base, cat: 'comp', out: classMap[mc[1].toUpperCase()] + mc[2] + mc[3].toUpperCase() + '-' + mc[4], added: mc[1] + '-' + mc[2] + mc[3] });
        } else add({ ...base, cat: 'en_exo', out: ev });
        continue;
      }
      // 其他：標籤數量不同，或（有原版時）ExoAE 確定有改
      const d = tagDiff(ev, zv);
      if (exoChanged === true || d.length) {
        add({ ...base, cat: 'check', out: zv, added: d.join(' ') });
        continue;
      }
      add({ ...base, cat: 'zh', out: zv });
    }
    // 自動合併的項目做檢查：有問題的改列「問題項目」，其餘直接採用、不必審核
    for (const e of entries) {
      if (!AUTO_CATS.has(e.cat)) continue;
      const iss = validate(e);
      if (iss.length) { stats[e.cat]--; stats.issue = (stats.issue || 0) + 1; e.from = e.cat; e.cat = 'issue'; e.issues = iss; }
    }
    // 礦石：壓成「中文＋英文縮寫」，每行不超過指定字數
    if (opts.ore !== false) {
      const ores = oreProposals(entries, zhT, exo, gloss, opts.ore || {});
      for (const e of entries) {
        const o = ores.get(e.k); if (!o) continue;
        stats[e.cat]--; stats.ore = (stats.ore || 0) + 1;
        e.from = e.cat; e.cat = 'ore'; e.out = o.out; e.note = o.note; e.len = o.len;
      }
    }
    // 地點：英文在前、中文在後
    if (opts.loc !== false) {
      const L = locProposals(entries, exo, opts.loc || {});
      for (const e of entries) {
        const v = L.get(e.k); if (v === undefined || v === e.out) continue;
        stats[e.cat]--; stats.loc = (stats.loc || 0) + 1;
        e.from = e.cat; e.cat = 'loc'; e.out = v;
      }
    }
    // 飛船：中文在前、英文在後（中文名來自星際公民汉化组）
    if (opts.ships && opts.ships.size) {
      for (const e of entries) {
        if (!/^vehicle_Name/.test(e.k)) continue;
        const z = (opts.ships.get(e.k) || '').replace(/(?:\\n|\s)+$/, '').trim(), en = (e.out || '').replace(/(?:\\n|\s)+$/, '').trim();
        if (!z || !en || CJK.test(en) || z.includes(en)) continue;
        stats[e.cat]--; stats.ship = (stats.ship || 0) + 1;
        // 太長就把英文換到下一行（中文算 1 字、英文數字算半字）
        // 版面：wrap＝一律中文一行、英文下一行（預設）；auto＝超過 shipWrap 字才換行；line＝同一行
        const layout = opts.shipLayout || 'wrap', wrap = opts.shipWrap == null ? 15 : +opts.shipWrap;
        const one = z + ' ' + en, two = z + '\\n' + en;
        e.from = e.cat; e.cat = 'ship'; e.shipZh = z;
        e.out = layout === 'line' ? one : layout === 'auto' ? (wrap > 0 && textLen(one, true) > wrap ? two : one) : two;
      }
    }
    // 漢化固定不動、ExoAE 改版時：英文變了的條目，中文可能已過時
    if (opts.exoBase) {
      for (const e of entries) {
        // 輸出完全由 ExoAE 產生的分類（礦石、零件等級、未中文化）會自動跟著新版，不用列出
        if (['exo_only', 'issue', 'ore', 'comp', 'en_exo', 'same', 'force', 'fill', 'ship'].includes(e.cat) || (e.cat === 'loc' && e.from !== 'zh')) continue;
        const bv = opts.exoBase.map.get(e.k);
        if (bv === undefined || bv === e.exo) continue;
        stats[e.cat]--; stats.updated = (stats.updated || 0) + 1;
        e.from = e.from || e.cat; e.cat = 'updated'; e.oldExo = bv;
      }
    }
    const seen = new Set(exo.keys);
    for (const k of zhRaw.keys) if (!seen.has(k)) add({ k, cat: 'zh_only', out: zhT.get(k), zh: zhT.get(k), zhRaw: zhRaw.map.get(k), exo: undefined, van: van ? van.map.get(k) : undefined });
    return { entries, stats, gloss };
  }

  // 漢化包漏翻、由你補上的翻譯（鍵名 => 中文）
  const DEFAULT_FILL = `# 鍵名 => 中文。漢化包漏翻的條目在這裡補上，會直接取代輸出
mission_items_QTQuantumBeacon => 量子信標
mission_items_QTDistressBeacon => 遇險信標
mission_items_QTPatrolBeacon => 巡邏信標
mission_items_QTRendezvousBeacon => 會合信標
mission_items_BlackBox => 飛行記錄器
mission_items_Confidential => 機密
mission_items_Derelict => 殘骸
item_NameBRAA_groundturret_AA => AAT-34 砲台
item_NameBRAA_groundturret_AA_short => AAT-34 砲台
item_NameBEHR_groundturret_AA_S4 => AAT-21 砲台
item_NameBEHR_groundturret_AA_S4_short => AAT-21 砲台
`;

  // 一律用 ExoAE 原文、不翻譯的條目（一行一個鍵名）
  const DEFAULT_FORCE_EXO = `# 一行一個條目鍵名，這些一律用 ExoAE 原文
comm_HackedConnection
`;

  // ---------- 地點 ----------
  // 漢化包沒翻的地點（多半是 Stanton、Pyro）可以在這裡補中文
  const DEFAULT_LOC_ZH = `# 英文 => 中文。只用在漢化包沒翻、仍是英文的地點名稱
Stanton => 斯坦頓
Hurston => 赫斯頓
Crusader => 十字軍
ArcCorp => 弧光星
microTech => 微科星
Yela => 耶拉
Calliope => 卡利俄佩
Clio => 克利俄
Euterpe => 歐忒耳佩
Area18 => 18區
Levski => 列夫斯基
Grim HEX => 六角灣
Pyro => 派羅
Ruin Station => 廢墟空間站
Stanton Gateway => 斯坦頓星門
Pyro Gateway => 派羅星門
Port Olisar => 奧麗莎空間站
`;
  function locProposals(entries, exo, o) {
    // 星系名稱：有 XxxSystem / Xxx_Star 鍵，或值是「Xxx System」
    const sys = new Set();
    for (const k of exo.keys) {
      let m = /^([A-Z][a-z]+)(?:System|_Star)$/.exec(k); if (m) sys.add(m[1]);
      if (/^[A-Z][a-z]+$/.test(k) && / System$/.test(exo.map.get(k) || '')) sys.add(k);
    }
    const re = new RegExp('^(?:(?:' + [...sys].join('|') + ')(?:\\d+[a-z]?)?(?:_[A-Za-z0-9_\\-]+)?|RR_[A-Za-z0-9_]+|ui_pregame_port_\\w+_name)(?:,P)?$');
    const bad = /desc|Desc|DESC|_Org|Courier|_from|[Hh]int|[Mm]ission|[Tt]itle|_obj|[Ss]hop|RepUI|Allies|Enemies/;
    const dict = new Map();
    for (const d of parseDict(o.zhNames != null ? o.zhNames : DEFAULT_LOC_ZH)) if (d.s) dict.set(d.s, d.to);
    const out = new Map();
    const esc = x => x.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    for (const e of entries) {
      if (!re.test(e.k) || bad.test(e.k) || e.exo == null || e.zh == null) continue;
      if (!['zh', 'same', 'en_exo'].includes(e.cat)) continue;
      const en = e.exo.trim();
      if (!en || en.length > 60 || /[~%<:\\]/.test(en) || CJK.test(en)) continue;
      let zh = e.zh.trim();
      if (CJK.test(zh)) {
        zh = zh.replace(new RegExp('\\s*[（(]\\s*' + esc(en) + '\\s*[)）]\\s*', 'g'), '').trim();   // 艾若(Aero) → 艾若
        if (!zh || zh.includes(en)) continue;
        out.set(e.k, en + ' ' + zh);
      } else {
        const z = dict.get(en);
        if (z) out.set(e.k, en + ' ' + z);
      }
    }
    return out;
  }

  // ---------- 零件等級 ----------
  const DEFAULT_CLASS_MAP = { CIV: '民用', MIL: '軍規', IND: '工業', STL: '隱形', STH: '隱形', CMP: '競速' };
  // MIL-1C "Bracer"  /  Mil/1/B Delphi RS-10
  const COMP_RE = /^(CIV|MIL|IND|STL|STH|CMP)[-\/]S?(\d+)\/?([A-Da-d])\s+"?(.+?)"?\s*$/i;

  // ---------- 礦石 ----------
  // 預設中文名（遊戲漢化沒有的才用；標 * 的是建議譯名，請自行確認）
  const DEFAULT_ORE_ZH = `# 英文 => 中文（漢化包裡已有的礦石會自動抓，不用寫在這裡；寫了會優先用這裡的）
Aluminum => 鋁
Diamond => 鑽石
Sadaryx => 薩達瑞晶
Savrilium => 薩維里
Carinite => 科力晶
Saldynium => 薩迪尼金
Beradom => 貝拉頓石*
Feynmaline => 費曼電氣石*
Glacosite => 格拉科石*
Jaclium => 賈克金屬*
Janalite => 賈納石*`;
  const DEFAULT_ORE_ABBR = `# 英文 => 縮寫。寫在這裡的一定照用（就算超過字數）；沒寫的自動取能互相分辨的開頭 3–6 個字母
Quantainium => Quan
Tungsten => Tung
Savrilium => Sav
Saldynium => Saldy
Sadaryx => Sadaryx
Carinite => Carinite
`;
  const ORE_SUFFIX = { ore: '礦', raw: '原', pure: '純' };
  const ALIAS = { aluminium: 'aluminum', savrillium: 'savrilium' };

  function textLen(s, half) {
    let n = 0;
    for (const ch of s) n += half && ch.charCodeAt(0) < 0x2e80 ? 0.5 : 1;
    return n;
  }
  function oreProposals(entries, zhT, exo, gloss, o) {
    const limit = o.limit || 14, half = !!o.halfWidth;
    const zhOver = new Map(), abbrOver = new Map();
    for (const d of parseDict(o.zhNames != null ? o.zhNames : DEFAULT_ORE_ZH)) if (d.s) zhOver.set(d.s.toLowerCase(), d.to);
    for (const d of parseDict(o.abbr != null ? o.abbr : DEFAULT_ORE_ABBR)) if (d.s) abbrOver.set(d.s.toLowerCase(), d.to);

    // 1. 找出所有礦石條目
    const scan = [], comm = [];
    for (const k of exo.keys) {
      let m = /^mineabletype_primary_([a-z]+)$/i.exec(k);
      if (m) { scan.push({ k, base: ALIAS[m[1].toLowerCase()] || m[1].toLowerCase() }); continue; }
      m = /^items_commodities_(raw_)?([a-z]+?)(_ore|_raw|_pure)?(,P)?$/i.exec(k);
      if (m) comm.push({ k, base: m[2].toLowerCase(), suf: (m[3] || (m[1] ? '_raw' : '')).slice(1).toLowerCase() });
    }
    const scanBases = new Set(scan.map(x => x.base));
    const oreBases = new Set(scanBases);
    for (const c of comm) if (c.suf) oreBases.add(ALIAS[c.base] || c.base);
    for (const c of comm) c.base = ALIAS[c.base] || c.base;
    const commOres = comm.filter(c => oreBases.has(c.base));

    // 2. 英文全名、中文名
    const en = new Map(), zh = new Map(), sug = new Set();
    const cap = b => b[0].toUpperCase() + b.slice(1);
    for (const c of commOres) if (!c.suf) { const v = exo.map.get(c.k); if (v && !/[()]/.test(v)) en.set(c.base, v.trim()); }
    for (const s of scan) {
      const v = (exo.map.get(s.k) || '').replace(/\s*\(.*$/, '').trim(); if (v && !en.has(s.base)) en.set(s.base, v);
      const z = zhT.get(s.k); if (z && CJK.test(z)) zh.set(s.base, z.trim());
    }
    for (const b of oreBases) {
      if (!en.has(b)) en.set(b, cap(b));
      const over = zhOver.get(en.get(b).toLowerCase()) || zhOver.get(b);
      if (over) { zh.set(b, over.replace(/\*$/, '')); if (/\*$/.test(over)) sug.add(b); }
      else if (!zh.has(b)) { const g = gloss.get(en.get(b)); if (g) zh.set(b, g); }
    }
    // 3. 縮寫：取能互相分辨的最短開頭（至少 3 個字母）
    const names = [...oreBases].map(b => en.get(b).replace(/[^A-Za-z]/g, ''));
    const minLen = new Map();
    for (const b of oreBases) {
      const w = en.get(b).replace(/[^A-Za-z]/g, '');
      let L = 3;
      while (L < w.length && names.some(x => x !== w && x.slice(0, L).toLowerCase() === w.slice(0, L).toLowerCase())) L++;
      minLen.set(b, Math.min(L, 6));
    }
    const abbrs = b => {
      const full = en.get(b), w = full.replace(/[^A-Za-z]/g, '');
      const o2 = abbrOver.get(full.toLowerCase()) || abbrOver.get(b);
      if (o2) return [o2];               // 你指定的寫法：一定照用
      const list = [full];
      for (let n = 6; n >= minLen.get(b); n--) if (n < w.length) list.push(w.slice(0, n));
      return [...new Set(list)];
    };
    const fit = s => s.split('\\n').every(l => textLen(l, half) <= limit);
    const out = new Map();
    const put = (k, s, b, extra) => out.set(k, { out: s, len: s.split('\\n').map(l => textLen(l, half)), note: (sug.has(b) ? '中文是建議譯名，請確認' : '') + (extra || '') });

    // 每種礦石的所有條目（掃描、原礦、精煉…）用同一個縮寫：取能讓「全部條目」都放得下的最長寫法
    const jobs = new Map(); // base -> [{k, one: a=>str, two: a=>str, plain}]
    const job = (b, j) => { if (!jobs.has(b)) jobs.set(b, []); jobs.get(b).push(j); };
    for (const s of scan) {
      const b = s.base, c = zh.get(b) || en.get(b), rs = ((exo.map.get(s.k) || '').match(/RS\s*([\d,]+)/) || [])[1];
      job(b, { k: s.k, one: a => `${c} ${a}${rs ? ' ' + rs : ''}`, two: a => rs ? `${c} ${a}\\nRS ${rs}` : `${c}\\n${a}`, plain: rs ? `${c}\\nRS ${rs}` : c });
    }
    for (const x of commOres) {
      const b = x.base, c0 = zh.get(b) || en.get(b), c = x.suf ? `${c0}(${ORE_SUFFIX[x.suf] || x.suf})` : c0;
      job(b, { k: x.k, one: a => `${c} ${a}`, two: a => `${c}\\n${a}`, plain: c });
    }
    for (const [b, js] of jobs) {
      const cands = abbrs(b);
      if (cands.length === 1 && (abbrOver.get(en.get(b).toLowerCase()) || abbrOver.get(b))) {
        for (const j of js) put(j.k, j.one(cands[0]), b, fit(j.one(cands[0])) ? '' : '超過字數（照你指定的縮寫）');
        continue;
      }
      let pick = cands.find(a => js.every(j => fit(j.one(a))));
      if (pick) { for (const j of js) put(j.k, j.one(pick), b, !zh.has(b) ? '找不到中文名' : ''); continue; }
      // 一行放不下：英文縮到最短仍超過 → 改兩行（中文優先）
      pick = cands[cands.length - 1];
      for (const j of js) {
        const t = fit(j.one(pick)) ? j.one(pick) : fit(j.two(pick)) ? j.two(pick) : j.plain;
        put(j.k, t, b, (!zh.has(b) ? '找不到中文名 ' : '') + (t === j.one(pick) ? '' : '已改成兩行'));
      }
    }
    return out;
  }

  function buildIni(entries, pick, eol) {
    const lines = new Array(entries.length);
    for (let i = 0; i < entries.length; i++) {
      const e = entries[i];
      lines[i] = e.k + '=' + String(pick(e)).replace(/\r?\n/g, '\\n');
    }
    return '\ufeff' + lines.join(eol || '\n') + (eol || '\n');
  }

  const api = { parseIni, merge, buildIni, parseDict, applyDict, textLen, CJK, DEFAULT_TAIL_DICT, DEFAULT_TERM_DICT, DEFAULT_ORE_ZH, DEFAULT_ORE_ABBR, DEFAULT_CLASS_MAP, DEFAULT_FORCE_EXO, DEFAULT_LOC_ZH, DEFAULT_FILL };
  if (typeof module !== 'undefined' && module.exports) module.exports = api; else root.SCMerge = api;
})(this);
