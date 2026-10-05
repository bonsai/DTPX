# DTPX

**DTP / 印刷 / 出版 / デジタル出力の世界モデル。**

DTPXは「文字や絵をページに置く」ためのテンプレートではない。
編集者が考えた意味を、**点・線・面・色・インク・版・紙・光・コード・PDF**へ変換し、
最終的に「読める／見える／刷れる／配れる」体験へ到達させるためのオントロジである。

## World

```
意味
 ↓
素材
 ↓
ページ
 ↓
図形・文字・画像
 ↓
座標・dot・色
 ↓
Colorant / Ink / Separation
 ↓
PDL (PostScript / PDF)
 ↓
RIP / Ghostscript
 ↓
Plate / Printer
 ↓
Paper / Object
 ↓
Eye / Camera / Screen
 ↓
Experience
```

デジタル側には別の枝がある。

```
Meaning → HTML / CSS → Browser → Screen
Meaning → TeX → PDF → Screen / Print
Drawing → DXF → CAD / Plotter / Print
```

重要なのは「ファイル形式」を中心にしないこと。
中心にあるのは **表現（representation）と媒体（medium）と体験（experience）の変換**である。

## Core ontology

- **Meaning** — 編集・設計される意味
- **Content** — 文章、写真、図、表、注釈
- **Document** — 内容を構造化したもの
- **Page** — 空間として構成された文書
- **Graphic** — 点・線・面・文字の幾何表現
- **Color** — 見た目／色空間／色値
- **Colorant** — 実際に出力へ寄与する着色材
- **Ink** — 印刷機で使われる物質
- **Separation** — 色を版・インク単位へ分解したもの
- **Dot** — 連続量を印刷可能な網点・画素へ離散化したもの
- **PDL** — Page Description Language
- **RIP** — ページ記述を出力デバイス向けデータへ解釈・変換する境界
- **Plate** — 版
- **Printer** — 出力機
- **Paper** — 物理媒体
- **Screen** — デジタル表示媒体
- **Experience** — 読む、見る、触る、持つ、配るという最終体験

## Color is not Ink

RGB / CMYK / Spot Color / DeviceN / ICC / Ink は同じ概念ではない。

色は「見え方」を表し、Colorantは出力に関わる着色要素、Inkは物理的な材料である。
Ghostscriptの色管理でも、RGB/CMYK等のsource/destination、ICC profile、DeviceN、spot color、output intentなどが別レイヤとして扱われる。 citeturn0search0turn0search4

## PDL as boundary

PostScriptとPDFは単なる「保存形式」ではなく、ページを記述し、出力へ渡すための重要な境界である。
GhostscriptはPostScript/PDFのインタプリタとグラフィックス処理を提供し、pdfwriteやps2writeなどの出力デバイスを持つ。 citeturn0search1turn0search3

## DXF

DXFは図面世界との境界。
DTPXでは「本の図版」と「CAD図面」を同一視しないが、

```
Geometry → Drawing → Page → Output
```

という変換関係を共有する。

## Agent boundary

DTPXは万能エージェントではない。

| Agent / Domain | Boundary |
|---|---|
| editor-agent | 企画・編集・意味・構成 |
| make-book | 本を商品として企画・制作・出版・販売 |
| GenkoyoshiWriter | 原稿入力・原稿用紙体験 |
| DTPX | ページ、組版、色、出力、印刷、デジタル表現 |
| printer / print-shop agent | 印刷機、紙、インク、工程、見積、製造 |
| CAD / DXF agent | 図面・幾何・設計情報 |
| web agent | HTML/CSS/ブラウザ体験 |
| PDF/TeX agent | 文書組版・PDF生成 |
| asset agent | 写真・画像・フォント等の素材管理 |

**境界を越えるときは、意味を勝手に変更せず契約を渡す。**

## Experience loop

```
編集する
  ↓
組む
  ↓
見る
  ↓
Proofする
  ↓
変換する
  ↓
刷る / Renderする
  ↓
触る・読む
  ↓
観察する
  ↓
修正する
```

DTPXのゴールは「PDFができた」ではなく、**媒体上で意図した体験が成立した**ことである。

## Toilet-paper typesetting

トイレットペーパーを長尺の読書媒体として扱う。本文は媒体仕様から独立したJSONLを正本とし、横書き日本語・縦書き日本語・横書き英語へ同じ内容を再組版できる。

- media: 114mm × 30m, 2-ply, flexography, 1-color body print
- printable pattern: 75mm × 860mm, 20mm gap
- recommended prototype body: 9pt / 14pt line height
- `specs/toilet-paper.jsonl`: media and layout contracts
- `examples/interview-10000.jsonl`: 10,000-character interview structure example

Pipeline: `interview.jsonl → layout JSONL → 860mm patterns → SVG → PDF/EPS`
