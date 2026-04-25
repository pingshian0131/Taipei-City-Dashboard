---
name: page
description: 當使用者要建立新頁面（View）、新增路由、或建立包含組件與地圖的獨立頁面時觸發此技能。關鍵字：新頁面、新 View、addRoute、新增路由、page、頁面
---

# 建立新頁面

建立新頁面是一個多步驟任務。請依照以下順序執行，並在每個步驟確認相關規範。

## 開始前

向使用者確認以下資訊（如尚未提供）：

1. **頁面名稱與路由路徑**（例：`Dispatch1999View`、`/dispatch-1999`）
2. **Layout 類型**：
   - 左側組件 + 右側地圖（同 MapView）
   - 純組件 grid（同 DashboardView）
   - 其他自定義 layout
3. **資料來源**：使用既有 contentStore dashboard 流程，或自行呼叫 API
4. **組件清單**：每個組件的圖表類型、是否有地圖
5. **額外功能**：城市切換、收藏、定時更新、彈跳視窗等
6. **視覺對標**：這個頁面要跟**哪個既有頁面**長得像？（MapView / DashboardView / 其他？）
   **必答**，讓 Step 0 知道要讀哪個檔案。

## Step 0：觀摩既有實作（MANDATORY，不可跳過）

寫任何一行新 Vue 之前，**必須**讀使用者在 Q6 指定的對標頁面，並讀以下檔案（依頁面類型）：

| 頁面類型 | Step 0 必讀 |
|---|---|
| 含地圖（MapView 式）| `src/views/MapView.vue`、`src/components/map/MapContainer.vue`、`src/store/mapStore.js` 的 `initializeMapBox` / `addMapLayer` |
| 純組件 grid（DashboardView 式）| `src/views/DashboardView.vue`、`src/dashboardComponent/DashboardComponent.vue` props 清單 |
| 有圖表 | 一個既有 `src/dashboardComponent/components/*.vue`（對應 Q4 的圖表類型） |

讀完後，**在對話中回報以下對照**（不可省略）：

- 既有頁面的 UI 骨架（SideBar / SettingsBar / 卡片殼 / 地圖殼分別是哪些元件）
- 你打算**重用**哪些、**自幹**哪些
- 任何**明確偏離既有樣式**的點（例：「不用 DashboardComponent 殼，自己寫卡片 → 視覺會有差」）

使用者確認對照表後，才進 Step 1。

### 重用優先原則

- **預設**：重用既有殼與子組件（`DashboardComponent`、`MapContainer`、`dashboardComponent/components/*`、`ComponentTag`、`dialogStore`、`mapStore`）
- **Opt-out 條件**：使用者明示「獨立 demo 頁」「不要用共用元件」才自幹
- 即使 opt-out，CSS 變數（`--color-*`、`--font-*`）與 `rules/code-style.md` 的命名／順序規則**仍強制遵守**

## 步驟一：建立 View 元件

在 `src/views/` 下建立 Vue 元件。

- 命名規則：參照 .claude/rules/code-style.md 的「Vue 元件」段落
- 程式碼結構：參照 .claude/rules/code-style.md 的「Vue 元件」結構順序
- 樣式規範：參照 .claude/rules/uiux.md 的系統顏色與間距變量
- 詳細示範：參照 .claude/skills/page/reference/page-creation-guide.md

### 含地圖的頁面必備元素

- 引入 `MapContainer` 元件
- 引入 `mapStore`，處理 toggle / filter / fly 事件
- 用 `computed` 將組件分為 hasMap / noMap 兩組
- 實作 `handleToggle()`、`shouldDisable()` 函式

### 純組件 grid 頁面必備元素

- 使用 CSS grid 排版，參照 DashboardView.vue 的 media query 斷點
- 不需要引入 MapContainer 和 mapStore

## 步驟二：註冊路由

在 `src/router/index.js` 的 `routes` 陣列新增路由定義。

## 步驟三：設定路由守衛

在 `src/router/index.js` 中修改對應的 `router.beforeEach`：

- 內容載入守衛（約第 164 行）：加入新路由的資料載入邏輯
- 地圖清除守衛：含地圖的頁面不應清除 mapStore
- 行動裝置守衛：決定新頁面是否允許行動裝置存取
- 權限守衛：決定新頁面是否需要登入

## 步驟四：配置 App.vue layout

在 `src/App.vue` 的 template 中，為新頁面設定 SideBar、SettingsBar 等 layout 結構。

## 步驟五：加入圖表組件

如果頁面包含圖表：
- 圖表資料格式：讀取 .claude/skills/chart/reference/chart-data.md
- 圖表類型與設定：讀取 .claude/skills/chart/reference/chart-type.md
- 圖表元件結構：讀取 .claude/skills/chart/SKILL.md

## 步驟六：加入地圖圖層

如果頁面包含地圖：
- 地圖資料格式：讀取 .claude/skills/map/reference/map-data.md
- 地圖類型與設定：讀取 .claude/skills/map/reference/map-type.md
- 地圖篩選功能：讀取 .claude/skills/map/reference/map-filter.md
- 底圖與圖層基礎：讀取 .claude/skills/map/SKILL.md

## 步驟七：加入彈跳視窗（如需要）

如果頁面需要彈跳視窗：
- 讀取 .claude/skills/dialog/SKILL.md
- 在 dialogStore 註冊新彈跳視窗
- 彈跳視窗 Vue 元件放在觸發元素旁邊，不要重複放置

## 完成後檢查

- [ ] Vue 元件命名為 PascalCase 且至少兩個英文字
- [ ] CSS root class 與元件名一致（全小寫無空格）
- [ ] CSS 屬性順序符合 code-style.md
- [ ] 使用專案定義的 CSS 變量（--color-*、--font-*）
- [ ] 路由守衛正確處理資料載入與地圖清除
- [ ] 不產生 console.log（ESLint no-console 規則）
- [ ] 執行 `npm run lint` 確認無錯誤
