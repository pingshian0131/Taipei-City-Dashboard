# Claude Code 設定完整性審查

> 審查日期：2026-04-19
> 審查對象：`/Users/ray/Documents/black/Taipei-City-Dashboard/` 根目錄的 `CLAUDE.md` 與 `.claude/` 目錄
> 審查視角：前端工程師只負責 `Taipei-City-Dashboard-FE/` 子專案的協作情境
> 結論：**結構蠻完整，但有 2 個 🔴 critical 會直接誤動作、2 個 🟠 major 內部矛盾**。建議修完再協作。

---

## 目錄

- [📁 現況盤點](#現況盤點)
- [🔴 Critical 問題](#critical-問題)
- [🟠 Major 問題](#major-問題)
- [🟡 Minor 問題與補強建議](#minor-問題與補強建議)
- [✅ 做得好的部分](#做得好的部分)
- [📋 修正優先順序](#修正優先順序)
- [📝 修正範本](#修正範本)

---

## 現況盤點

### CLAUDE.md 內容（16 行）

| 區塊 | 內容 |
|---|---|
| Language | 強制繁體中文 ✓ |
| Project Overview | **引用 `ProjectOverview.md`（檔案不存在）** ❌ |
| 問題輸出 | **要求輸出到 `/docs`（目錄不存在，本審查順便建立）** ❌ |

### .claude/ 目錄結構

```
.claude/
├── agent/                      ⚠️ 應是 agents/（複數）
│   └── code-review.md
├── rules/
│   ├── code-style.md          ✓ 完整
│   └── uiux.md                ✓ 完整
└── skills/
    ├── chart/                 ✓
    ├── dialog/                ✓
    ├── map/                   ✓
    ├── page/                  ✓
    └── skill-creator/         ✓（元 skill，含 eval 工具）
```

### Taipei-City-Dashboard-FE/ 技術棧推導

從 `package.json` 推導：
- Vue 3.4（`<script setup>` Composition API）
- Pinia 2 state management
- Vue Router 4
- Vite 5（build & dev）
- SCSS（sass 1.70）
- Apexcharts 3.45 + vue3-apexcharts
- Mapbox GL 3.1 + deck.gl 9 + threebox
- axios、dayjs、lodash.debounce、@vueuse/core、turf、material-icons
- ESLint 9（flat config `eslint.config.js`）
- **非 TypeScript**（`main.js` 不是 `.ts`）
- **無 TailwindCSS**
- **無測試依賴**（vitest / jest / playwright 都沒裝）

### src/ 結構（FE 工作區）

```
Taipei-City-Dashboard-FE/src/
├── App.vue
├── main.js
├── assets/
├── components/
├── dashboardComponent/        （chart skill 提及的目錄）
├── directives/
├── router/
├── store/                     （Pinia store）
└── views/                     （page skill 建立的目標）
```

---

## 🔴 Critical 問題

### C1. `ProjectOverview.md` 不存在但 CLAUDE.md 強制引用

**位置**：[CLAUDE.md:11](../CLAUDE.md#L11)

```
ProjectOverview.md contains the project overview. Please read it first.
```

**問題**：
- 整個 repo `find -iname "ProjectOverview*"` 結果為空
- Claude 收到 CLAUDE.md 後會依指示先嘗試讀取 → 失敗 → 可能做出錯誤假設
- 浪費 token，每次任務都會多一次失敗 tool call

**修法二擇一**：
1. 建立 `ProjectOverview.md`（建議放在根目錄，內容見文末範本）
2. 若暫無內容，改寫這行為：「專案概覽見 README.md」並確認 README 有足夠 context

### C2. `/docs` 不存在但 CLAUDE.md 要求輸出到該目錄

**位置**：[CLAUDE.md:15](../CLAUDE.md#L15)

```
每次詢問類似架構或文件改善的內容時，把結果輸出到/docs裡，並新增一個md檔案，不要直接修改原有的檔案。
```

**問題**：
- `/docs` 不存在（本審查已順便建立 `docs/` 與 `docs/code-review-reports/`）
- `agent/code-review.md:135` 也寫「docs/code-review-reports/」作為報告輸出位置
- 第一次 Claude 會嘗試寫入失敗

**修法**：本審查已修（建立空目錄 + .gitkeep）。若要更完善，在 `docs/README.md` 標注子目錄用途：
```
docs/
├── code-review-reports/    (agent/code-review.md 輸出)
├── claude-setup-review.md  (本檔)
└── <其他改善建議文件>.md
```

---

## 🟠 Major 問題

### M1. `agent/code-review.md` 的標準與實際技術棧矛盾

**位置**：[.claude/agent/code-review.md:36–42](../.claude/agent/code-review.md#L36)

**矛盾 1：TailwindCSS vs SCSS**

```
- [ ] 使用 TailwindCSS classes，無 hardcoded colors
```

但：
- `package.json` 沒有 tailwindcss 依賴
- `rules/uiux.md` 明寫「本專案都使用純 css 或 scss 進行樣式設計」
- `rules/code-style.md` 的 CSS 範例都是 SCSS class + `var(--color-xxx)`

**矛盾 2：TypeScript vs JavaScript**

```
- [ ] 優先使用 `interface` 而非 `type`
```

`interface` / `type` 是 TypeScript 語法，但專案用 JavaScript：
- `main.js`（非 `.ts`）
- `package.json` 沒有 typescript 依賴
- `vite.config.js` 非 `.ts`

**矛盾 3：測試要求與專案現況**

```
🧪 測試 (Testing)
- [ ] 新功能是否有對應測試
```

但 `package.json` **沒有任何測試框架**（vitest／jest／playwright 皆無）。

**後果**：code-review agent 會對所有 PR 誤報大量問題，審查報告失信度。

**修法**：重寫 `agent/code-review.md` Step 3 檢查項目，對齊實際技術棧：

```markdown
#### 📝 程式碼品質 (Code Quality)
- [ ] 使用 Vue 3 Composition API + `<script setup>`
- [ ] Vue 元件命名為 PascalCase（至少兩個英文字）
- [ ] 函式名以動詞開頭 + camelCase
- [ ] 使用 SCSS + CSS variables（`var(--color-xxx)`），無 hardcoded hex
- [ ] root class 與檔案名相同但全小寫（參考 rules/code-style.md）
- [ ] CSS 屬性順序依 rules/code-style.md（dimensions → display → position → margin/padding → border → background → font → animation → transition → other）
- [ ] 使用 named functions / named exports；不使用 `var`
- [ ] 無 console.log、debugger、未使用的 imports

#### 🧪 測試（目前專案無測試框架）
- 若 PR 新增測試框架，需先在 docs/ 提架構決議
- 暫跳過此段

#### 🎨 專案風格規則
- 對照 rules/code-style.md 與 rules/uiux.md 逐項檢查
```

### M2. CLAUDE.md 未明示 FE 工作邊界

**問題**：這是 monorepo：
```
Taipei-City-Dashboard/
├── Taipei-City-Dashboard-FE/    ← 前端工程師負責
├── Taipei-City-Dashboard-BE/
├── Taipei-City-Dashboard-DE/
├── db-sample-data/
├── docker/
└── helm-chart/
```

但 CLAUDE.md **完全沒有說明前端工程師只在 FE 子專案工作**。Claude 可能：
- 誤動 BE 的 Go 程式碼
- 誤改 docker-compose 或 helm chart
- 在根目錄建立不該在根目錄的檔案

**修法**：在 CLAUDE.md 補上「工作邊界」段落（範本見文末）。

### M3. `.claude/agent/` 命名可能不會自動載入

**問題**：Claude Code 的標準 subagents 目錄是 `.claude/agents/`（複數），這個專案是 `.claude/agent/`（單數）。

**驗證**：
- 若 `code-review` agent 實際從未被觸發成功，確認原因就是目錄名
- 可先執行 `/agents` 看列表是否有 code-review

**修法**：重新命名 `agent/` → `agents/`，並檢查 agent file 的 frontmatter（`name:`、`description:`、`tools:`、`model:`）是否正確。

---

## 🟡 Minor 問題與補強建議

### N1. `rules/*.md` 沒在 CLAUDE.md 明確引用

`.claude/rules/` 底下的檔案 Claude Code 不一定會自動載入。建議在 CLAUDE.md 明確寫：

```
## 開發規則（MUST read）

- 程式碼風格：.claude/rules/code-style.md
- UI / UX 設計原則：.claude/rules/uiux.md
```

### N2. 缺技術棧與開發指令速查

CLAUDE.md 沒提：
- 技術棧（Vue 3 / Pinia / Mapbox / Apexcharts / deck.gl / SCSS / Vite）
- 開發指令（`npm run dev` / `npm run build` / `npm run lint`）
- `src/` 目錄組織

每次 Claude 都要自己探索 package.json + 目錄才能開工，浪費 token。

### N3. 缺 API 串接約定

FE 會打 BE API，但 CLAUDE.md 與 rules 都沒提：
- base URL（env 變數名）
- 認證方式（token？cookie？）
- FE 端 axios interceptor 的慣例

這會讓新開 API 串接的 PR 很容易偏離既有模式。

### N4. 缺 Git / PR 工作流

例如：
- 分支命名規則（`feature/xxx`、`fix/xxx`）
- Commit message 規範（中文？英文？Conventional Commits？）
- PR 標題格式、標籤規範
- code-review 何時觸發

### N5. `code-review.md` 輸出路徑與 CLAUDE.md 不一致

- CLAUDE.md：`把結果輸出到/docs裡`
- code-review.md:135：`docs/code-review-reports/`

後者是前者的子目錄，**實際不衝突**，但建議 CLAUDE.md 也明寫分類：

```
## 問題輸出

- 架構／文件改善：docs/<topic>.md
- Code review 報告：docs/code-review-reports/<PR>.md
- （其他類型）：docs/<category>/<filename>.md
```

---

## ✅ 做得好的部分

| 項目 | 評價 |
|---|---|
| `rules/code-style.md` | **完整**：Vue 元件骨架、函式／變數命名、CSS class 命名、CSS 屬性撰寫順序都有範例 |
| `rules/uiux.md` | **完整**：設計哲學（簡約）、顏色 CSS 變數、字體層級、圖示字體、局部樣式邊界 |
| `skills/page` | 觸發描述清楚（「新頁面、addRoute、新 View」）|
| `skills/chart` | 觸發描述清楚（dashboardComponent 資料夾、ApexCharts）|
| `skills/dialog` | 觸發描述清楚（dialogStore、DialogContainer、Teleport）|
| `skills/map` | 觸發描述清楚（mapStore、mapConfig、Mapbox）|
| `skills/skill-creator` | 有 eval 工具、viewer、benchmark scripts 等完整 meta-skill 基礎設施 |
| `agent/code-review.md` 流程 | Step 0–4 與輸出格式設計完整，問題在檢查項目標準（M1）|

---

## 📋 修正優先順序

| 優先 | 問題 | 工作量 | 影響 |
|:-:|---|:-:|---|
| 🥇 | **C1** 建立 ProjectOverview.md 或修掉引用 | 10 分鐘 | 每次 Claude 啟動都有效 |
| 🥇 | **C2** 建立 /docs（本審查已修） | — | ✓ |
| 🥈 | **M1** 改寫 code-review.md 檢查標準 | 20 分鐘 | code-review agent 變可信 |
| 🥈 | **M2** CLAUDE.md 補 FE 工作邊界 | 10 分鐘 | 避免誤動 BE/DE |
| 🥉 | **M3** 確認 agent/→agents/ 重新命名 | 5 分鐘 | code-review agent 可實際觸發 |
| 4 | **N1–N5** 補 CLAUDE.md 缺失 context | 30 分鐘 | 日常效率提升 |

---

## 📝 修正範本

### 建議的新版 CLAUDE.md 骨架

```markdown
# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Language

YOU MUST respond in 繁體中文 (zh-TW). NEVER use 簡體中文 (zh-CN).

## 專案結構（monorepo）

本 repo 包含三個子專案，目前協作者只負責前端：

- **`Taipei-City-Dashboard-FE/`** ← **你只在這個子目錄工作**
- `Taipei-City-Dashboard-BE/`（Go 後端，不要動）
- `Taipei-City-Dashboard-DE/`（資料工程，不要動）
- `docker/`、`helm-chart/`、`db-sample-data/`（不要動）

若需跨專案改動，請先在 `docs/` 提架構決議並告知使用者。

## 技術棧（FE）

- Vue 3.4 + Composition API（`<script setup>`）
- Pinia 2 / Vue Router 4 / Vite 5
- SCSS（pure，**無 TailwindCSS**）
- Apexcharts 3.45 + vue3-apexcharts
- Mapbox GL 3 + deck.gl 9 + threebox
- ESLint 9（flat config `eslint.config.js`）
- axios / dayjs / lodash.debounce / @vueuse/core / turf / material-icons
- **非 TypeScript（main.js）**
- **無測試框架**（若要新增需先提決議）

## 開發指令（在 Taipei-City-Dashboard-FE/ 底下）

- `npm run dev` — Vite dev server
- `npm run build` — lint 後 build（production）
- `npm run build:test` — test 模式 build
- `npm run lint` — ESLint auto-fix
- `npm run preview` — Vite preview

## src/ 目錄（FE）

- `views/` — 頁面層（新頁面用 page skill）
- `components/` — 共用 UI 元件
- `dashboardComponent/` — 儀表板內的圖表組件（圖表用 chart skill）
- `store/` — Pinia store
- `router/` — Vue Router
- `directives/` — 自定義 directive
- `assets/` — 靜態資源

## 開發規則（MUST read）

- 程式碼風格：[.claude/rules/code-style.md](.claude/rules/code-style.md)
- UI / UX 設計原則：[.claude/rules/uiux.md](.claude/rules/uiux.md)

## Skills（依使用情境自動觸發）

- page：建新頁面、新增路由
- chart：新增／修改 dashboardComponent 圖表
- dialog：建新彈跳視窗
- map：Mapbox 底圖、圖層、mapStore

## 問題輸出

每次詢問架構或文件改善類內容時，輸出到：

- 架構／文件改善：`docs/<topic>.md`
- Code review 報告：`docs/code-review-reports/<PR>.md`

**不要直接修改原有檔案**。
```

### 建議的新版 code-review.md Step 3

```markdown
### Step 3: 執行審查檢查項目

#### 📝 程式碼品質 (Code Quality)

- [ ] 使用 Vue 3 Composition API + `<script setup>`
- [ ] Vue 元件為 PascalCase（至少兩個英文字）
- [ ] 函式名以動詞開頭 + camelCase
- [ ] 使用 SCSS + CSS variables（`var(--color-xxx)`），無 hardcoded hex
- [ ] root class 與 Vue 檔名相同但全小寫
- [ ] CSS 屬性順序遵循 rules/code-style.md
- [ ] 使用 named functions / named exports；不使用 `var`
- [ ] 無 console.log、debugger、未使用 imports

#### 🎨 專案風格規則 (Project Style Rules)

對每個變更檔案，用 Grep 驗證是否符合：
- rules/code-style.md（Vue 元件骨架、命名、CSS 順序）
- rules/uiux.md（顏色／字體層級 CSS 變數、元件命名一致性）

#### 🧪 測試

專案目前無測試框架。若 PR 新增測試框架相關檔案（vitest.config.js 等），需要求先在 docs/ 提架構決議。

#### 🗂 架構

- [ ] 新增檔案是否在 FE src/ 的正確子目錄（views/components/dashboardComponent/store/router）
- [ ] 有無誤改 BE/DE 子專案或 docker/helm 設定
```

---

## 總結

現有設定的**骨架是好的**（rules、skills 都用心寫過），但有兩個直接讓 Claude 誤動作的 critical 問題（引用不存在的檔案 + 輸出到不存在的目錄），以及 code-review agent 檢查標準與專案實際技術棧完全對不上的 major 問題。

**建議修正順序**：C1 + C2 先做（已順手把 C2 做掉），再處理 M1 讓 code-review 可信，最後補 M2 / M3 / N 類提升日常效率。修完後這個 `.claude/` 設定對前端工程師協作來說會很扎實。
