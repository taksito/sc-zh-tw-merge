"""產生 src/live_missing_zh_s.txt：ExoAE 漏掉、但遊戲正式版會用到的鍵（簡體原文，工具會轉繁體）。

比對 MrKraken StarStrings 正式版（master）與 ExoAE：ExoAE 沒有的鍵，取星際公民汉化组「4.10.3 LIVE 汉化临时版」(main) 的中文。
用法：python3 tools/make_live_missing.py starstrings_master.ini hanhua_main.ini
"""
import sys
from pathlib import Path
root = Path(__file__).resolve().parent.parent
def parse(p):
    d = {}
    for l in open(p, encoding='utf-8-sig'):
        l = l.rstrip('\r\n')
        if '=' in l and not l.startswith(';'):
            k, v = l.split('=', 1); d[k] = v
    return d
mk, zh = parse(sys.argv[1]), parse(sys.argv[2])
exo, half = parse(root / 'data/exocomp.ini'), parse(root / 'data/chinese.ini')
out = [f'{k}={zh[k].strip()}' for k in mk if k not in exo and k not in half and zh.get(k, '').strip()]
head = ['# ExoAE 漏掉、遊戲正式版會用到的鍵（簡體，工具會轉繁體並強制補上）',
        '# 鍵名來源：MrKraken StarStrings 正式版；中文：星際公民汉化组 4.10.3 LIVE 汉化临时版 V0.1（2026-11-15）']
MANUAL = '# ---- 手動補充：官方英文檔也沒有、遊戲卻會用到的鍵（重新產生時保留）----'
def manual(path):
    t = path.read_text(encoding='utf-8') if path.exists() else ''
    return t[t.index(MANUAL):].rstrip('\n').split('\n') if MANUAL in t else [MANUAL]
zh_manual, en_manual = manual(root / 'src/live_missing_zh_s.txt'), manual(root / 'src/live_missing_en.txt')
(root / 'src/live_missing_zh_s.txt').write_text('\n'.join(head + out + zh_manual) + '\n', encoding='utf-8')
keys = [l.split('=', 1)[0] for l in out]
en = ['# ExoAE 漏掉的鍵的英文原文（取自 MrKraken StarStrings 正式版）；ExoAE 補上後就不會用到'] + [f'{k}={mk[k].strip()}' for k in keys]
(root / 'src/live_missing_en.txt').write_text('\n'.join(en + en_manual) + '\n', encoding='utf-8')
ships = [f'{k}={zh[k].strip()}' for k in keys if k.lower().startswith('vehicle_name')]
print('ships', ships)
print(len(out), 'keys')
