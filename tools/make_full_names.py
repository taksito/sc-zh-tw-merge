"""產生 src/full_names_zh.txt：全繁中版用的「全漢化」補充中文。

半漢化包（data/chinese.ini）故意把地面名稱、載具、物品等留成英文。全繁中版要全部中文，
所以從星際公民汉化组「简中全汉化」(StarCitizenToolBox/LocalizationData main 分支) 抽出：
  半漢化仍是英文、全漢化已有中文 的條目（星圖地名、太空站、拉格朗日點、勢力、人名、任務物品…）。
物品名稱（item_names_full_zh.txt）、船名（ship_names_zh.txt）、任務地點（mission_locations_zh.txt）已有專用表，這裡不重複收。
變數記號（~xxx()、%s、<EM>、\\n 數量）跟目前英文不一致的條目跳過，避免遊戲顯示壞掉。

用法：python3 tools/make_full_names.py <全汉化 global.ini> [版本說明]
"""
import re, sys
from pathlib import Path

root = Path(__file__).resolve().parent.parent
def parse(p):
    d = {}
    for l in open(p, encoding='utf-8-sig'):
        l = l.rstrip('\r\n')
        if '=' in l and not l.startswith(';'):
            k, v = l.split('=', 1); d[k.replace(',P', '')] = v
    return d
def table_keys(p):
    return {l.split('=', 1)[0] for l in open(p, encoding='utf-8-sig') if '=' in l and not l.startswith('#')}

cjk = lambda s: bool(re.search(r'[㐀-鿿]', s or ''))
TOK = re.compile(r'~\w+\([^)]*\)|%\w|<\/?EM\d?>|\\n')
toks = lambda s: sorted(TOK.findall(s or ''))
SKIP = re.compile(r'^(vehicle_name|mission_location|hud_)', re.I)

full = parse(sys.argv[1]); ver = sys.argv[2] if len(sys.argv) > 2 else ''
half, exo = parse(root / 'data/chinese.ini'), parse(root / 'data/exocomp.ini')
have = table_keys(root / 'src/item_names_full_zh.txt') | table_keys(root / 'src/mission_locations_zh.txt')
out, skipped = [], []
for k in sorted(exo):
    en, f = exo[k], full.get(k)
    if SKIP.match(k) or k in have or cjk(half.get(k)) or not cjk(f) or cjk(en): continue
    # \n 只比數量；~mission() 這類變數要一模一樣
    # 括號要成對（漢化組有一條被截斷成「[已关闭」），[CLOSED]→[已关闭] 這種翻譯是正常的
    if [t for t in toks(en) if t != '\\n'] != [t for t in toks(f) if t != '\\n'] or en.count('\\n') != f.count('\\n') or f.count('[') != f.count(']'):
        skipped.append(k); continue
    out.append(f'{k}={f.strip()}')
# 手動補：漢化組那條被截斷（「[已关闭」少了右括號）
MANUAL = {'Stanton4a_UGF_001': '舒宾加工设施 SPMC-3 [已关闭]'}
have_out = {l.split('=', 1)[0] for l in out}
out += [f'{k}={v}' for k, v in MANUAL.items() if k not in have_out and k in exo]
head = ['# 全繁中版用的「全漢化」補充中文（簡體原文，工具會轉成繁體）',
        '# 來源：星際公民汉化组「简中全汉化」https://github.com/StarCitizenToolBox/LocalizationData （main 分支' + (f'，{ver}' if ver else '') + '）',
        '# 只收半漢化包仍是英文、全漢化已有中文的條目（星圖地名、太空站、勢力、人名…）；全繁中版只在輸出仍是英文時才套用',
        '# 由 tools/make_full_names.py 產生，請勿手動修改']
(root / 'src/full_names_zh.txt').write_text('\n'.join(head + out) + '\n', encoding='utf-8')
print(len(out), 'entries;', len(skipped), 'skipped (token mismatch):', skipped[:15])
