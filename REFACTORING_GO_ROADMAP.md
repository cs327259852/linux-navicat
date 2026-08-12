# 🚀 LinuxNavicat Go + Wails 渐进式重构路线图 (Daily Roadmap)

本项目正在从原有的 Java Swing 架构升级重构为现代化的 **Go + Wails v2/v3 + Vue 3 / React + TypeScript** 极轻量桌面客户端。

---

## 🎯 重构目标

1. **极致轻量**：安装包/可执行文件体积 < 15MB，运行内存占用 < 50MB。
2. **现代美观 UI**：Web 前端渲染，支持极简 Dark/Light 主题切换、流畅平滑动画。
3. **高并发高并发 SQL 引擎**：依托 Go 语言原生的并发与数据库 Driver，SQL 执行极速且不阻塞界面。
4. **多数据库扩展能力**：支持 MySQL、PostgreSQL、SQLite、Redis 统一接入。

---

## 📅 每日优化落地计划 (Daily Progressive Milestones)

### 🚩 Day 1: 环境搭建与工程脚手架初始化
- [ ] **环境准备**：安装 Go (>= 1.20) 与 Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)。
- [ ] **项目脚手架初始化**：
  ```bash
  wails init -n LinuxNavicat-Go -t vue-ts
  ```
- [ ] **构建验证**：运行 `wails dev`，确保 Go + TS 前后端通信与窗口展示正常。

### 🚩 Day 2: Go 后端核心连接模块与加密存储
- [ ] **连接数据结构**：定义 `ConnectionConfig` 结构体（Host, Port, User, Password, Database）。
- [ ] **数据库驱动引入**：引入 `go-sql-driver/mysql`。
- [ ] **连接测试服务**：实现 Go 方法 `TestConnection(config ConnectionConfig) (bool, string)`。
- [ ] **本地安全存储**：配置 JSON 持久化与简单的 AES 加密处理密码存储。

### 🚩 Day 3: 前端侧边栏与连接管理模态框
- [ ] **侧边栏组件**：使用 Tailwind CSS + Lucide Icons 构建连接树与菜单面板。
- [ ] **新建连接弹窗**：支持配置填写、测试连接按钮与保存提示。
- [ ] **连接状态指示**：在线/离线状态图标展示与错误 Toast 提示。

### 🚩 Day 4: 数据库与数据表树形结构导航
- [ ] **后端元数据导出**：
  - `GetDatabases(connId string) ([]string, error)`
  - `GetTables(connId string, dbName string) ([]TableInfo, error)`
- [ ] **前端虚拟树组件**：展开数据库节点、异步懒加载所有数据表列表。

### 🚩 Day 5: 现代化 SQL 编辑器集成 (Monaco Editor / CodeMirror 6)
- [ ] **前端集成编辑**：引入 `@monaco-editor/react` 或 `CodeMirror 6`。
- [ ] **SQL 语法高亮与补全**：配置 MySQL 关键字自动补全、括号匹配。
- [ ] **快捷键支持**：按 `Ctrl+Enter` / `Cmd+Enter` 快速执行选中或当前 SQL。

### 🚩 Day 6: 高性能数据网格 (DataGrid) 与分页控制
- [ ] **后端查询接口**：`ExecuteQuery(connId, sql, page, pageSize)`。
- [ ] **前端数据表格**：展示动态列名与行数据，支持单元格文本复制。
- [ ] **分页条工具栏**：首页、上一页、下一页、末页与数据总条数计算。
- [ ] **查看 DDL**：一键弹出当前表的 `SHOW CREATE TABLE` 建表语句。

### 🚩 Day 7: 数据导出与暗黑/亮色主题实时切换
- [ ] **数据导出功能**：前端/后端导出查询结果为 CSV、JSON 或 Excel (.xlsx)。
- [ ] **主题系统**：集成 Tailwind `dark` 模式，支持一键切换黑夜/白天主题。

### 🚩 Day 8: 跨平台 GitHub CI/CD 自动化打包与 Release
- [ ] **配置 GitHub Actions**：编写 `.github/workflows/wails-build.yml`。
- [ ] **多平台构建**：自动生成 Linux (Deb/AppImage)、macOS (DMG)、Windows (EXE)。
- [ ] **自动 Release**：在 GitHub 上发布对应 Tag 的二进制供用户下载试用。

---

## 🛠️ 准备开始 Day 1

准备好开始第一天（Day 1）的初始化与环境配置了吗？只需安装 Go 环境即可协同推进！
