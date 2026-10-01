# TeX 結びつけ — tex_to_dtpx 契約の実装仕様
#
# 対象: bonsai/sotsusei/thesis（総合演習 提出物を pandoc + tectonic で PDF にする）
# 検証: 2026-10-01 / pandoc 3.12 + tectonic 0.17 (XeTeX) on Windows
#
# このファイルは「TeX を DTPX の PDL 境界へ接続する」を実務仕様として書いたもの。
# model/contracts.yml の tex_to_dtpx skeleton を flesh out した実装仕様である。

contract: tex_to_dtpx
input:
  - md_or_tex      # 意味（構造化された原稿）
  - template       # 版面設計（documentclass / geometry / フォント）
  - fonts          # CJK + 欧文
  - assets         # 画像
output:
  - pdf            # PDL 境界の出力
  - log            # 診断（どのページが組まれたか）

---

## 1. 境界の位置

DTPX において TeX は「意味 → 幾何 → PDL」の**一実装**であり、
PDF/TeX agent と editor-agent の境界を横切る。以下のように位置づける。

```
  editorial agent        tex_to_dtpx 契約        DTPX 側
  ──────────────         ─────────────────        ──────────
  Markdown (意味)  ──▶  md → tex      ──▶
  LaTeX (構造)
                       pandoc  (意味→構造)
                       template (版面→構造)
                       tectonic (構造→PDF)
                                                 PDF ──▶ RIP / Screen

  返るのは PDF と log のみ。意味は持ち戻さない（DTPX の境界則）。
```

**重要な主張**: TeX を「黒箱 PDF 変換器」として扱うのは誤り。
LaTeX の版面（geometry / baselineskip / 段落組）は**幾何の生成**であり、
DTPX の `composition → graphics` に相当する。ログと .tex を残すのは
この段を観測可能にするため。

---

## 2. 契約違反（実測）

pandoc の既定テンプレートを前提にして_template を書くと、次の
6 点で**契約が破れる**。いずれも「pandoc が出力する形」に対する
DTPX 側の受皿の不足であり、DTPX 側で持つべきではない。

### V1. longtable / \LTcaptype

- **pandoc の出力**:
  ```
  {\def\LTcaptype{none} % do not increment counter
  \begin{longtable}[]{p{...}p{...}}
  ...
  \end{longtable}}
  ```
- **受信側が壊れる理由**: longtable.sty は `\begin{longtable}` 内で
  `\addtocounter{LTcaptype}` を呼ぶ。pandoc が先に `\def\LTcaptype{none}`
  を書いているため、`\setcounter{none}` にフォールし
  `LaTeX Error: No counter 'none' defined.` で停止する。
- **longtable 環境を再定義する逃げ**は不可。booktabs の `\noalign` が
  `Misplaced \noalign` / `\TX@get@body has an extra }` で壊れる。
- **受容**: SPEC を含む複数行を扱える**行ベースの後処理**で
  `\begin{table}` + `\begin{tabular}` へ書き換える。

### V2. SPEC 内の \real

- pandoc は列幅に `p{(\linewidth - 8\tabcolsep) * \real{0.2000}}` と書く。
- `\real` は longtable の拡張で、tabular には無い。
- **受容**: 後処理で `\real{x}` → `x` に落とす。

### V3. \NormalTok / Shaded（ハイライト系）

- pandoc は必ず `\begin{Shaded}\begin{Highlighting}[..]\NormalTok{..}` を出力する。
- 独自 template では両方未定義 → `Environment Shaded undefined` /
  `Undefined control sequence \NormalTok`。
- **罠**: `\DefineVerbatimEnvironment{Highlighting}{Verbatim}` を使うと
  `\NormalTok` が verbatim 内で展開されず印字される。**不可**。
- **受容**: `Shaded` / `Highlighting` は空撕え、文字列トークンは
  素通し（`#1`）で定義する。色分けは pandoc の既定 handler に任せる。

### V4. \tightlist

- pandoc は必ず `\tightlist` を出力する。独自 template では未定義。
- **受容**: `\providecommand{\tightlist}{...}` で用意する。

### V5. biblio-style 未設定

- metadata に `bibliography` が無い場合、既定テンプレートは
  `\bibliographystyle{plain}` を無条件に出し、`plain.bst` 解決で失敗する。
- **受容**: `$if(bibliography)$` で囲む。

### V6. title / subtitle

- `\subtitle` は article.cls に無い。
- `\maketitle` を再定義すると `\@author` 展開が vertical mode で
  `You can't use \spacefactor in vertical mode` で落ちる。
- **受容**: `\title{\shortstack[l]{...$title$...$subtitle$...}}` に畳む。
  再定義しない。

---

## 3. フォント（Windows + tectonic）

### 3.1 Fontconfig が無い

tectonic の XeTeX は Windows のシステムフォントを探索できない。
`Fontconfig error: Cannot load default config file` が常に出るが、
**警告であり致命的ではない**。TeX Gyre Termes のような
*tectonic バンドル内* のファイル名なら解決できる。

→ Windows の BIZ-UDMincho.ttc 等のシステムフォントは
  `The font "..." cannot be found` になる。

### 3.2 使用可能なフォント

| 用途 | 名前 | 備考 |
|---|---|---|
| 欧文 serif | `texgyretermes-*.otf` | tectonic 同梱。4 種 |
| 欧文 sans | `texgyreheros-*.otf` | 同梱 |
| 欧文 mono | `texgyrecursor-regular.otf` | 同梱 |
| CJK serif | `ipaexm.ttf` | JIS 漢字のみ |
| CJK sans | `ipaexg.ttf` | 同上 |
| CJK fallback | `FandolSong-Regular.otf` | 簡体字用 |

### 3.3 CJK の二つの制約

1. **太字が無い**: `ipaexb.ttf` / `ipaexg.ttf` は tectonic に無い。
   → `BoldFeatures={FakeBold=1.8}` で擬似太字を作る。
2. **簡体字が無い**: `进 / 图 / 语 / 总` 等が欠落する。
   → `\setCJKfallbackfamilyfont{\CJKrmdefault}{FandolSong-Regular.otf}`

日本語のみなら fallback は不要。**日本語と簡体字が混在する原稿では必須**。

### 3.4 罫線文字

`┌ ─ │ ┐` 等の box-drawing は `texgyrecursor` に無い。
`\texttt{}`（Highlighting 内）だと `Missing character U+250C` になる。
框線を使うなら CJK フォント側（ipaexg.ttf）に字体がある。

---

## 4. 観測可能性

### 4.1 契約の戻り値

`output` に `log` を明示した理由：

- TeX は**失敗が本文中に紛れる**。行番号すら出ないことがある。
- DTPX の `feedback` エッジ（experience → design_change）に
  「どの段で落ちたか」を渡すために必要がある。

### 4.2 観測すべき点

```
note: Running TeX ...
note: Writing `out.pdf` (N KiB)
warning: ... Missing character: There is no X (U+XXXX) in font [...]
error:   file.tex:NNN: <message>
```

- `Writing` が出れば PDF は生成された（**終了コードより Reliable**）。
- `Missing character` は font の話で、**致命的ではない**。
  `error:` 行だけを失敗判定に使う。
- tectonic は `texput.tex` という一時名で走らせるため、
  終了コードだけ見ると原因ファイルが分からない。`--keep-logs` で
  ログを残すこと。

---

## 5. 受入基準（実装の check）

```
[ ] markdown 1 ファイル → PDF 1 ファイルが生成される
[ ] CJK が豆腐（□）にならない
[ ] 表が 1 ページに収まり、罫線が引かれる
[ ] code block が等幅で表示される
[ ] \LTcaptype / \real / \NormalTok 由来の error が 0 件
[ ] Missing character warning が 0 件（框線を使うなら代替字体を指定）
[ ] log に出力ファイル名とページ数が記録される
```

---

## 6. まだ決まっていないこと

| # | 未決 | 影響 |
|---|---|---|
| 1 | 1 ページを跨ぐ表|longtable を本当に使うか。跨がなければ V1 の後処理で足りる |
| 2 | 日本語フォントの Georgia | IPAex は旧字体。HaranoAji を tectonic に追加する方法が未検証 |
| 3 | 参考文献 | bibtex を通すか。手書きで `\begin{thebibliography}` にするか |
| 4 | 図の挿入 | graphicx は入れたが、座標指定は未検証 |

---

## 関連

- contracts: `model/contracts.yml` の `tex_to_dtpx`
- 実装: `bonsai/sotsusei` の `thesis/template.tex` / `thesis/build.ps1`
- 既存アセット: PDF/TeX agent（README の Agent boundary 表）
