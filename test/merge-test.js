// 用法：node test/merge-test.js chinese.ini exocomp.ini [english.ini]
// 不做簡轉繁，只檢查合併邏輯與分類數量
const M = require('../src/engine.js'), fs = require('fs');
const [zhP, exP, vanP] = process.argv.slice(2);
if (!zhP || !exP) { console.error('用法：node test/merge-test.js chinese.ini exocomp.ini [english.ini]'); process.exit(1); }
const zh = M.parseIni(fs.readFileSync(zhP, 'utf8')), ex = M.parseIni(fs.readFileSync(exP, 'utf8'));
const van = vanP ? M.parseIni(fs.readFileSync(vanP, 'utf8')) : null;
const r = M.merge(zh, new Map(zh.map), ex, van, {});
console.log(r.stats);
for (const e of r.entries.filter(e => e.cat === 'issue')) console.log('問題：', e.k, e.issues.join('；'));
