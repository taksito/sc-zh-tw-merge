"""產生 src/item_names_full_zh.txt：全繁中版用的物品中文名稱。

來源：星際公民汉化组「简中全汉化」(StarCitizenToolBox/LocalizationData)
  - 4.10.2 正式版全汉化（release/4.10.2_v2_20261014）為主
  - 英文名稱在 4.10.3 有改的鍵，改用較新的 5.0 PTU 全汉化（main）
只收「我們用的半漢化包仍是英文」的物品／零件名稱。

用法：python3 tools/make_item_names.py full_4102.ini full_main.ini
"""
import re, sys
from pathlib import Path

root = Path(__file__).resolve().parent.parent
def parse(p):
    d = {}
    for l in open(p, encoding='utf-8-sig'):
        l = l.rstrip('\r\n')
        if '=' in l and not l.startswith(';'):
            k, v = l.split('=', 1); d[k] = v
    return d

cjk = lambda s: bool(re.search(r'[㐀-鿿]', s or ''))
NAME_RE = re.compile(r'^(item_name|items_|vehicle_name)', re.I)
f2, f5 = parse(sys.argv[1]), parse(sys.argv[2])
half = parse(root / 'data/chinese.ini')
exo, base = parse(root / 'data/exocomp.ini'), parse(root / 'data/exocomp.base.ini')
out = []
for k in sorted(set(exo) | set(half)):
    if not NAME_RE.match(k) or cjk(half.get(k)): continue
    changed = k in exo and exo.get(k) != base.get(k)
    v = f5.get(k) if changed and cjk(f5.get(k)) else f2.get(k) if cjk(f2.get(k)) else f5.get(k) if cjk(f5.get(k)) else None
    if v: out.append(f'{k}={v.strip()}')
# 手動補：漢化組也沒有、遊戲卻會用到的名稱（簡體）
MANUAL = {
    'item_Nameksar_lmg_energy_01': '裁决 轻机枪',
    'item_Nameksar_lmg_energy_01_short': '裁决 轻机枪',
    'item_Nameksar_lmg_energy_01_store01': '裁决 轻机枪 衍化',
    'item_Nameksar_lmg_energy_01_store01_short': '裁决 轻机枪 衍化',
}
have = {l.split('=', 1)[0] for l in out}
out += [f'{k}={v}' for k, v in MANUAL.items() if k not in have]
head = ['# 全繁中版用的物品中文名稱（簡體原文，工具會轉成繁體）',
        '# 來源：星際公民汉化组「简中全汉化」https://github.com/StarCitizenToolBox/LocalizationData',
        '#   4.10.2(PU)_CN_V2（release/4.10.2_v2_20261014）為主；4.10.3 英文有改的鍵改用 5.0(PTU)_CN_V1（main，2026-11-13）',
        '# 只收半漢化包仍是英文的物品／零件名稱；格式 鍵名=中文']
(root / 'src/item_names_full_zh.txt').write_text('\n'.join(head + out) + '\n', encoding='utf-8')
print(len(out), 'names')
