# =============================================================================
# HTML → PDF — 最短経路と、その代償
# =============================================================================
#
# 問い: HTML から PDF へ変換する方が、TeX より楽ではないですか。
#
# 答え: **はい、短期的には圧倒的に楽。** ただし条件つき。
#
# =============================================================================

summary:
  verdict: yes_with_conditions
  one_line: >-
    HTML→PDF は最も短いコードパスだが、
    「組版の結果」を受け取るのであって「組版の決定」を受け取るわけではない。

# -----------------------------------------------------------------------------
# なぜ楽か
# -----------------------------------------------------------------------------
why_easier:

  code_length:
    browser_route: ~10 行（headless Chrome / Playwright の print-to-PDF）
    tex_route: >-
      pandoc の既定テンプレートを前提にすると 6 箇所の契約違反が出る。
      実測: longtable/\LTcaptype, \real, \NormalTok/Shaded,
      \tightlist, biblio-style, \subtitle
    ratio: およそ 1 : 20

  what_browser_gives_free:
    - CSS による段落組版・行分割・日本語禁則処理
    - web font による日本語フォント解決
    - 表のautomatic 配置
    - flexbox/grid によるレイアウト
    - page-break 制御（@page, break-inside）
    - 成熟した hyphenation/kerning
    # いずれも TeX では自前 or テンプレートで対策していた領域

  font_handling:
    - IPAex の旧字体・太字欠落・簡体字欠落这些问题が**全部消える**
    - Noto Sans JP / Noto Serif JP を指定する
    - fake bold 不用
    - fallback 不用

  installation:
    browser_route: 追加インストールなし（Chrome は既存）
    tex_route: >-
      tectonic Asked（約 20MB）+ フォントバンドル取得。
      さらに Windows では Fontconfig が無くシステムフォントを探索できない。

# -----------------------------------------------------------------------------
# ただし — 代償
# -----------------------------------------------------------------------------
the_cost:

  # --- 1. 意味を復元できない ------------------------------------------------
  loses_semantic_structure:
    detail: >-
      HTML→PDF は「表示されたもの」を返す。元の Markdown は消える。
      見出し・引用・リストという**構造**は、
      PDF 上ではレイアウトの結果としてしか現れない。
    dtopx_impact: >-
      DTPX の editorial → composition の境界が失われる。
      「なぜこの改行が起きたか」を後から説明できない。
      feedback エッジ（experience → editorial_change）が機能しない。

  # --- 2. 再現性が保証されない ------------------------------------------------
  reproducibility:
    browser_route:
      depends_on: [chrome_version, installed_fonts, gpu_rasterization, os]
      # 同じ HTML でも Chrome を更新すると PDF が変わる
    tex_route:
      depends_on: [engine_version, pinned_revision]
      # tectonic は revision を pin できる
    verdict: >-
      卒制の提出物として**再ビルド可能性**が必要ななら TeX 側が有利。
      学会提出は査読者が再現する。

  # --- 3. 印刷用の制御が弱い --------------------------------------------------
  print_control:
    html_paged_media:
      supports: [@page size, @page margin, break-before/after, orphans, widows]
      weak: [bleed, crop_marks, color_bar, registration_marks, trap]
    tex_route:
      supports: [geometry, eso-pic での自由配置, graphicx, 完全に任意]
    # 実用成果物なら HTML で十分。
    # 印刷・製本列为 则 TeX（または PostScript）が必要的。

  # --- 4. 色分離が難しい ------------------------------------------------------
  color_separation:
    html_css_color_mode: print-color-adjust: exact
    problem: >-
      ブラウザは device RGB で組む。CMYK 分離は後段の RIP に委ねる。
      PostScript の `setcmykcolor` のような**明示的な着色材指定**が無い。
    dtopx_impact: >-
      ontology の separation → plate → ink の経路が、
      HTML 側では**发生时而不是設計時**に確定する。
      DTPX の設計としては正しいが、
      「どのインクを使うか」を編集時に決められない。

  # --- 5. ファイル ─────────────────────────────────────────────────────
  file_size:
    browser_route: >-
      web font を埋め込むため PDF が大きくなりがち。
      subset されずに全文字が落ちることもある。
    tex_route: 同程度（font embedding あり）
    verdict: どちらも最適化が必要

# -----------------------------------------------------------------------------
# 使い分け
# -----------------------------------------------------------------------------
recommendation:

  html_to_pdf_when:
    - 印刷に出さない（画面・要件定義・レビュー用）
    - 版面設計が CSS で表現できる範囲
    - 短納期
    - 日本語が絡む（TeX のフォント制約を避けたい）
    - 再現性が不要

  tex_when:
    - 再ビルド可能性が要件（学会提出・査読）
    - 版面（geometry / baseline / 章番号）を精密に制御したい
    - 数式が重い
    - 生成する PDF を diff したい（テステSuiteで stable に）
    - 印刷・製本列入 输出

  hybrid_when:
    - 主力は HTML→PDF（速さ）
    - 特定の成果物だけ TeX（再現性）
    - → 今回の総合演習はこれ。HTML で围着ULES、提出 PDF は Chrome で焼く。

# -----------------------------------------------------------------------------
# 実装（最短経路）
# -----------------------------------------------------------------------------
implementation:

  chrome_headless:
    command: >-
      chrome --headless --disable-gpu --print-to-pdf=out.pdf
      --no-pdf-header-footer --print-to-pdf-no-header
      file:///C:/path/to/doc.html
    pros: [追加依存なし, フォント解決が正しい]
    cons: [version 依存, 制御が限られる]

  playwright:
    pip: pip install playwright && playwright install chromium
    code: |
      from playwright.sync_api import sync_playwright

      HTML = """
      <html><head><meta charset="utf-8">
      <style>
        @page { size: A4; margin: 25mm; }
        body  { font-family: "Noto Serif JP", serif; line-height: 1.7; }
        @media print { h2 { break-after: avoid; } }
      </style></head><body>...</body></html>
      """

      with sync_playwright() as p:
          browser = p.chromium.launch()
          page = browser.new_page()
          page.set_content(HTML, wait_until="networkidle")
          page.pdf(path="out.pdf", format="A4",
                   margin={"top": "25mm", "bottom": "25mm",
                           "left": "25mm", "right": "25mm"},
                   print_background=True)
          browser.close()
    pros: [再現性可控（chromium を pin できる）, ヘッダ/フッタ制御, ページ番号]
    cons: [chromium のダウンロード (~150MB)]

  # -----------------------------------------------------------------------------
  # Markdown → HTML は素直
  # -----------------------------------------------------------------------------
  markdown_to_html:
    command: pandoc in.md -o out.html --standalone --css=paper.css
    note: >-
      LaTeX 変換（-t latex）で起きる 6 箇所の non-互換は、
      HTML 変換では**ゼロ**。pandoc の HTML writer は自己完結している。
    stylesheet_hint: |
      @page { size: A4; margin: 25mm; }
      body {
        font-family: "Noto Serif JP", "Yu Mincho", serif;
        line-height: 1.8;
        font-size: 11pt;
      }
      h1 { break-before: page; }
      h2 { break-after: avoid; }
      table { border-collapse: collapse; width: 100%; break-inside: avoid; }
      th, td { border: 1px solid #333; padding: 6px 10px; }
      code { font-family: "Noto Sans Mono", monospace; }

# -----------------------------------------------------------------------------
# 結論
# -----------------------------------------------------------------------------
conclusion: >-
  「楽か」で言えば Yes。ただし樂bedaく得来るのは
  「正しい PDF」であって「説明できる組版」ではない。
  DTPX の goal が「読める/見える/刷れる/配れる」体験へ到達することなら、
  HTML→PDF で十分。
  goal が「編集の意図を追跡できること」なら TeX の方が強い。
  設計としては両方を Lowering の実装として持ち、
  どちらを選ぶかは下流の要件で決めるのが妥当。
