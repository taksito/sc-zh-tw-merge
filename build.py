#!/usr/bin/env python3
"""把 src/app.html 與 src/engine.js 組成單一檔案的 index.html（GitHub Pages 直接用）。"""
import pathlib
root = pathlib.Path(__file__).parent
app = (root / 'src/app.html').read_text(encoding='utf-8')
eng = (root / 'src/engine.js').read_text(encoding='utf-8')
dec = (root / 'src/chinesedecode.txt').read_text(encoding='utf-8-sig')
head = ('<!doctype html>\n<html lang="zh-Hant">\n<head>\n<meta charset="utf-8">\n'
        '<meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover">\n'
        '<style>:root{color-scheme:light dark}body{margin:0}[hidden]{display:none!important}img{max-width:100%}</style>\n'
        '</head>\n<body>\n')
(root / 'index.html').write_text(head + app.replace('/*ENGINE*/', eng).replace('/*DECODE*/', dec) + '\n</body>\n</html>\n', encoding='utf-8')
print('index.html built')
