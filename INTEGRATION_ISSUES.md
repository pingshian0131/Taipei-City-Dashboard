# Taipei City Dashboard 整合問題清單

產出日期：2026-04-26
產出範圍：FE (`Taipei-City-Dashboard-FE`) + BE (`Taipei-City-Dashboard-BE`) + DE (`Taipei-City-Dashboard-DE`) + DB
分支：`ping-dev`（merge 進 `chester/dashboard-1`、`ray/dashboard-fe-1`）

---

## 🔴 嚴重：FE / BE 命名與 shape 不對齊

mrt-a11y 模組四個 endpoint 幾乎是 1:1 錯位：

| FE 期待 | BE 提供 | DB 表 | 狀態 |
|---|---|---|---|
| `/alert-count` | `/alert-count` ✅ | `mrtp_a11y_alert` | 已對齊（已拿掉 mock） |
| `/alert-by-line` | `/alert-by-line` ✅ | `mrtp_a11y_alert` | 已對齊（已拿掉 mock） |
| `/alert-by-type` | （無） | — | **FE 還在 mock**，BE 沒對應 |
| `/station-overview` | （無） | — | **FE 還在 mock**，BE 沒對應 |
| （無） | `/alert-trend-30d` | `mrtp_a11y_alert_history` | **BE 孤兒** |
| （無） | `/stations` | `mrtp_a11y_elevator` | **BE 孤兒** |

→ **雙方各做各的**，需要選一邊對齊：FE 改打 BE 既有 / BE 補出 FE 期待 / 兩邊改。

## 🟡 DB 表沒人用

| 表 | 行數 | BE | FE | 結論 |
|---|---|---|---|---|
| `mrtp_a11y_alert` | 1 (0 active) | ✅ | ✅ | UI 上是 0，因為 ETL 抓到的真實狀態就是無 active 公告 |
| `mrtp_a11y_alert_history` | 10 | ✅ trend-30d | ❌ | 完全沒上儀表板 |
| `mrtp_a11y_elevator` | 188 / 118 stations | ✅ stations | ❌ | 完全沒上儀表板 |

## 🟡 FE 殘留靜態資料

| 位置 | 內容 | scope |
|---|---|---|
| `Taipei-City-Dashboard-FE/mock/index.js` 還剩 2 條 entry | `alert-by-type`, `station-overview` | 中 |
| `Taipei-City-Dashboard-FE/mock/mrt-a11y/*.json` 4 個檔 | `alert-count.json`, `alert-by-line.json` 已 dead（mapping 已移除） | 清理 |
| `Taipei-City-Dashboard-FE/src/views/AccessibilityRouteView.vue:25-111` 寫死 | `slopeCounts`, `workCounts`, `slopeMockGeoJson`, `workMockGeoJson` — BE 跟 DB 都沒對應 | 大（需 ETL→DB→BE→FE 一條龍） |

## 🟡 Demo 資料缺乏

`mrtp_a11y_alert` 表只有 1 筆 `closed`（奇岩站已修復），導致即使串通了 BE：

- C1 顯示 `0`
- C2 空圖
- C4 (Map) 全綠

要看 UI 動起來必須塞測試 active 資料，**或等真實 ETL 抓到當天有公告才行**。

範例 SQL（會被下次 DAG 跑覆蓋掉）：

```sql
INSERT INTO mrtp_a11y_alert (line, station, publish_time, description, status, data_time) VALUES
  ('淡水信義線', '台北車站', NOW(), '電梯維修中',     'active', NOW()),
  ('板南線',     '忠孝復興', NOW(), '坡道暫停使用',   'active', NOW()),
  ('松山新店線', '西門',     NOW(), '月臺隙縫板施工', 'active', NOW());
```

## 🟢 本 session 已修復 / 改動

| 項目 | 狀態 |
|---|---|
| Airflow webserver port 衝突（8080 撞 FE） | 已改 `Taipei-City-Dashboard-DE/docker/develop/docker-compose.yaml` 8081:8080 |
| Vite docker 模式缺 `/api/v1` proxy | 已補 `Taipei-City-Dashboard-FE/vite.config.js` |
| Mock 攔截 `alert-count` / `alert-by-line` | 已從 `Taipei-City-Dashboard-FE/mock/index.js` 移除 |
| FE / BE Docker image rebuild | 已執行 |
| FE peer dep 衝突 (`@vitejs/plugin-basic-ssl@2.x` 跟 `vite@5`) | 用 `--legacy-peer-deps` 裝起來 — **未根治** |
| Nginx (port 80) 502 | FE 切回 HTTP 後可能自動好了，沒重新驗證 |

## 🟢 設計層面（次要）

- **FE 沒有 API service layer**：每個 view 直接 `axios.get(...)`，路徑寫死分散，難一次替換 base path / 加全域錯誤處理
- **BE `/alert-by-line` 跟 `/alert-trend-30d` 兩個 query 各自寫 SQL**：line 分組邏輯重複，可抽 helper
- **`AccessibilityRouteView` 註解寫「BE 將以行政區彙總回傳」「每 10 分鐘輪詢」**：需求已寫但 BE / DAG 沒實作 → 要決策是否真的做
- **`@vitejs/plugin-basic-ssl@2.x` 強拉了 vite 8 peer，但專案實際是 vite 5**：要不降 plugin 版本，要不升 vite，目前用 `--legacy-peer-deps` 跳過

## 已驗證的服務狀態（2026-04-26）

| 服務 | URL | 狀態 |
|---|---|---|
| FE (Vite) | http://localhost:8080（user 後續關了 basicSsl）| 200 |
| BE (Gin) | http://localhost:8088 | OK（root 404 正常） |
| Nginx | http://localhost | 待重驗 |
| pgAdmin | http://localhost:8889 | 302 → 登入 |
| Airflow UI | http://localhost:8081/airflow-sit/ | healthy |
| Postgres `dashboard` (postgres-data) | 容器內 5432 | running |
| Postgres `dashboardmanager` (postgres-manager) | localhost:5432 | running |
| Redis | 內部 6379 | running |
| Qdrant | 6333 (yaml 有定義) | **未啟動** |

---

## 建議優先順序

1. **塞測試資料到 `mrtp_a11y_alert`**（5 分鐘）— 讓 C1/C2 UI 看得出在動
2. **mrt-a11y 對齊收尾**（1-2 hr）— `alert-by-type` / `station-overview` / `alert-trend-30d` / `stations` 四個對齊一次處理乾淨；同時把 `mrtp_a11y_alert_history` + `mrtp_a11y_elevator` 兩個孤兒表上儀表板
3. **AccessibilityRouteView slope/work**（半天 ~ 1 天）— ETL + table + endpoint + FE，整套做
4. **FE service layer 抽出**（清債，可緩）
5. **`@vitejs/plugin-basic-ssl` 版本問題根治**（清債，可緩）

---

## 對齊方案備選（給 mrt-a11y 模組）

### 選項 A：BE 補出 FE 期待

| 新 BE endpoint | 從哪算 |
|---|---|
| `/alert-by-type` | 由 `mrtp_a11y_alert` 的 `description` 解析設施類型聚合（電梯 / 坡道 / 月臺隙縫板…），或要 ETL 多抽欄位 |
| `/station-overview` | 由 `mrtp_a11y_elevator` LEFT JOIN `mrtp_a11y_alert` 得每站異常 / 正常統計 |

優點：FE 不動。缺點：BE 要再加 endpoint，且 `alert-by-type` 的 facility_type 解析邏輯是新工作。

### 選項 B：FE 改打 BE 既有

| FE 改打 | 同時要做 |
|---|---|
| `/stations`（取代 `/station-overview`） | C4 元件改成地圖點位 marker（紅綠著色），不再是 overview 統計 |
| `/alert-trend-30d` 新增 C5 元件（取代 `/alert-by-type`） | 新增折線圖元件，30 天趨勢 |

優點：直接利用現有 BE。缺點：FE 重做兩個 component。

### 選項 C：兩邊都動

統一命名，例如：
- `/mrt/a11y/summary`（取代 alert-count + alert-by-type）
- `/mrt/a11y/by-line`（保留）
- `/mrt/a11y/trend?days=30`（取代 alert-trend-30d）
- `/mrt/a11y/stations`（保留）

優點：最乾淨。缺點：FE / BE 都要改。
