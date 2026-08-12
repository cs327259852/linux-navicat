<template>
  <div class="app-container" :class="{ 'light-theme': isLightTheme }">
    <!-- Sidebar -->
    <aside class="sidebar">
      <div class="brand">
        <span class="icon">🐬</span>
        <h2>LinuxNavicat</h2>
        <span class="version-tag">Go 0.2.0</span>
      </div>

      <div class="action-bar">
        <button class="btn btn-primary" @click="openNewConnectionModal">+ 新建连接</button>
      </div>

      <!-- Connection & Database Tree -->
      <div class="tree-section">
        <div class="tree-header">
          <span>数据库导航</span>
          <button class="btn-icon" title="刷新列表" @click="fetchConnections">🔄</button>
        </div>

        <ul class="tree-list">
          <li v-for="conn in connections" :key="conn.id" class="tree-node">
            <div class="node-header" :class="{ active: selectedConn?.id === conn.id }" @click="toggleConnection(conn)">
              <span class="tree-arrow" :class="{ expanded: conn.isExpanded }">▸</span>
              <span class="status-dot online"></span>
              <span class="node-label" :title="conn.name">{{ conn.name }}</span>
              <span class="node-sub">{{ conn.host }}:{{ conn.port }}</span>
              <button class="btn-delete" @click.stop="confirmDeleteConnection(conn)">✕</button>
            </div>

            <!-- Databases List -->
            <ul v-if="conn.isExpanded" class="tree-sub-list">
              <li v-for="dbName in conn.databases" :key="dbName" class="tree-node">
                <div class="node-header db-header" :class="{ active: selectedDb === dbName && selectedConn?.id === conn.id }" @click="toggleDatabase(conn, dbName)">
                  <span class="tree-arrow" :class="{ expanded: conn.expandedDbs?.[dbName] }">▸</span>
                  <span class="icon-db">📁</span>
                  <span class="node-label">{{ dbName }}</span>
                </div>

                <!-- Tables List -->
                <ul v-if="conn.expandedDbs?.[dbName]" class="tree-sub-list">
                  <li v-for="tableName in conn.dbTables?.[dbName]" :key="tableName" class="table-item" @click="selectTable(conn, dbName, tableName)">
                    <span class="icon-table">📄</span>
                    <span class="node-label table-label">{{ tableName }}</span>
                  </li>
                </ul>
              </li>
            </ul>
          </li>
        </ul>
      </div>
    </aside>

    <!-- Main Workspace Area -->
    <main class="main-content">
      <!-- Top Navigation Header -->
      <header class="top-header">
        <div class="tab-list">
          <div class="tab-item active">
            <span>SQL 查询标签页</span>
          </div>
        </div>

        <div class="top-controls">
          <button class="btn btn-theme" @click="isLightTheme = !isLightTheme">
            {{ isLightTheme ? '🌙 暗色模式' : '☀️ 亮色模式' }}
          </button>
          <button class="btn btn-run" :disabled="isExecuting" @click="runSQL">
            {{ isExecuting ? '⏳ 执行中...' : '▶ 执行 SQL (Ctrl+Enter)' }}
          </button>
        </div>
      </header>

      <!-- SQL Input Editor -->
      <div class="editor-pane">
        <textarea
          v-model="sqlCode"
          class="sql-textarea"
          placeholder="在此输入 SQL 查询语句... (例如: SELECT * FROM information_schema.tables;)"
          @keydown.ctrl.enter="runSQL"
          @keydown.meta.enter="runSQL"
        ></textarea>
      </div>

      <!-- Execution Status Bar & Actions -->
      <div class="status-bar">
        <div class="status-info">
          <span v-if="queryResult" :class="queryResult.success ? 'text-success' : 'text-danger'">
            {{ queryResult.message }}
          </span>
          <span v-if="queryResult?.executionTime != null" class="execution-time">
            (耗时: {{ queryResult.executionTime }} ms)
          </span>
        </div>

        <div v-if="queryResult?.rows?.length" class="export-actions">
          <button class="btn-export" @click="exportData('csv')">📥 导出 CSV</button>
          <button class="btn-export" @click="exportData('json')">📥 导出 JSON</button>
          <button v-if="currentTable" class="btn-export" @click="fetchTableDDL">🔍 查看 DDL</button>
        </div>
      </div>

      <!-- Result Table Grid -->
      <div class="result-pane">
        <div v-if="!queryResult" class="empty-state">
          <span>尚未执行 SQL 查询。请选择左侧数据表或输入 SQL 后点击【执行 SQL】</span>
        </div>

        <div v-else-if="!queryResult.success" class="error-box">
          <div class="error-title">❌ 执行发生错误:</div>
          <pre class="error-msg">{{ queryResult.message }}</pre>
        </div>

        <div v-else class="grid-wrapper">
          <table class="data-grid">
            <thead>
              <tr>
                <th class="col-index">#</th>
                <th v-for="col in queryResult.columns" :key="col">{{ col }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(row, rIdx) in paginatedRows" :key="rIdx">
                <td class="col-index">{{ (currentPage - 1) * pageSize + rIdx + 1 }}</td>
                <td v-for="(cell, cIdx) in row" :key="cIdx" class="cell-content" :title="cell">
                  {{ cell }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination Controls -->
        <div v-if="queryResult?.rows?.length" class="pagination-bar">
          <span class="page-info">共 {{ queryResult.rows.length }} 条记录 | 页码 {{ currentPage }} / {{ totalPages }}</span>
          <div class="page-btn-group">
            <button class="btn-page" :disabled="currentPage === 1" @click="currentPage = 1">首页</button>
            <button class="btn-page" :disabled="currentPage === 1" @click="currentPage--">上一页</button>
            <button class="btn-page" :disabled="currentPage === totalPages" @click="currentPage++">下一页</button>
            <button class="btn-page" :disabled="currentPage === totalPages" @click="currentPage = totalPages">末页</button>
          </div>
        </div>
      </div>
    </main>

    <!-- Modal: Create / Edit Connection -->
    <div v-if="showConnModal" class="modal-backdrop">
      <div class="modal-card">
        <div class="modal-header">
          <h3>连接 MySQL 数据库</h3>
          <button class="modal-close" @click="showConnModal = false">✕</button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label>连接名称 (Connection Name):</label>
            <input v-model="formConn.name" type="text" placeholder="例: 本地测试环境" />
          </div>

          <div class="form-row">
            <div class="form-group flex-2">
              <label>主机地址 (Host):</label>
              <input v-model="formConn.host" type="text" placeholder="127.0.0.1" />
            </div>
            <div class="form-group flex-1">
              <label>端口 (Port):</label>
              <input v-model.number="formConn.port" type="number" placeholder="3306" />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group flex-1">
              <label>用户名 (Username):</label>
              <input v-model="formConn.user" type="text" placeholder="root" />
            </div>
            <div class="form-group flex-1">
              <label>密码 (Password):</label>
              <input v-model="formConn.password" type="password" placeholder="••••••••" />
            </div>
          </div>

          <div class="form-group">
            <label>默认数据库 (Database Schema):</label>
            <input v-model="formConn.schema" type="text" placeholder="mysql" />
          </div>

          <div v-if="testResultMsg" class="test-feedback" :class="testResultSuccess ? 'text-success' : 'text-danger'">
            {{ testResultMsg }}
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" :disabled="isTesting" @click="testConnection">
            {{ isTesting ? '⏳ 测试中...' : '🔌 测试连接' }}
          </button>
          <div class="right-actions">
            <button class="btn btn-secondary" @click="showConnModal = false">取消</button>
            <button class="btn btn-primary" @click="saveConnection">保存并连接</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal: View Table DDL -->
    <div v-if="showDDLModal" class="modal-backdrop">
      <div class="modal-card modal-large">
        <div class="modal-header">
          <h3>数据表 DDL 建表语句 ({{ currentTable }})</h3>
          <button class="modal-close" @click="showDDLModal = false">✕</button>
        </div>
        <div class="modal-body">
          <pre class="ddl-code">{{ tableDDLContent }}</pre>
        </div>
        <div class="modal-footer">
          <button class="btn btn-primary" @click="showDDLModal = false">关闭</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

interface ConnectionConfig {
  id?: string
  name: string
  host: string
  port: number
  user: string
  password?: string
  schema: string
  isExpanded?: boolean
  databases?: string[]
  expandedDbs?: Record<string, boolean>
  dbTables?: Record<string, string[]>
}

interface QueryResult {
  success: boolean
  message: string
  columns?: string[]
  rows?: string[][]
  affectedRows?: number
  executionTime?: number
}

// State
const isLightTheme = ref(false)
const connections = ref<ConnectionConfig[]>([])
const selectedConn = ref<ConnectionConfig | null>(null)
const selectedDb = ref<string>('')
const currentTable = ref<string>('')

const sqlCode = ref<string>('SELECT * FROM information_schema.tables LIMIT 20;')
const queryResult = ref<QueryResult | null>(null)
const isExecuting = ref(false)

// Pagination
const currentPage = ref(1)
const pageSize = ref(50)

// Modals
const showConnModal = ref(false)
const isTesting = ref(false)
const testResultSuccess = ref(false)
const testResultMsg = ref('')

const showDDLModal = ref(false)
const tableDDLContent = ref('')

const formConn = ref<ConnectionConfig>({
  name: '本地 MySQL',
  host: '127.0.0.1',
  port: 3306,
  user: 'root',
  password: '',
  schema: 'mysql'
})

// Wails bindings helper
const wails = (window as any).go?.main?.App

const totalPages = computed(() => {
  if (!queryResult.value?.rows?.length) return 1
  return Math.ceil(queryResult.value.rows.length / pageSize.value)
})

const paginatedRows = computed(() => {
  if (!queryResult.value?.rows) return []
  const start = (currentPage.value - 1) * pageSize.value
  return queryResult.value.rows.slice(start, start + pageSize.value)
})

onMounted(() => {
  fetchConnections()
})

async function fetchConnections() {
  if (wails?.GetConnections) {
    try {
      const res = await wails.GetConnections()
      connections.value = (res || []).map((c: ConnectionConfig) => ({
        ...c,
        isExpanded: false,
        databases: [],
        expandedDbs: {},
        dbTables: {}
      }))
    } catch (e) {
      console.error('获取连接列表失败', e)
    }
  } else {
    // Mock fallback for browser testing
    connections.value = [
      { id: 'conn_1', name: 'Demo Local Database', host: '127.0.0.1', port: 3306, user: 'root', schema: 'mysql', isExpanded: false }
    ]
  }
}

function openNewConnectionModal() {
  formConn.value = {
    name: 'MySQL Server',
    host: '127.0.0.1',
    port: 3306,
    user: 'root',
    password: '',
    schema: 'mysql'
  }
  testResultMsg.value = ''
  showConnModal.value = true
}

async function testConnection() {
  isTesting.value = true
  testResultMsg.value = ''
  try {
    if (wails?.TestConnection) {
      const res = await wails.TestConnection(formConn.value)
      testResultSuccess.value = res.success
      testResultMsg.value = res.message
    } else {
      testResultSuccess.value = true
      testResultMsg.value = '测试连接成功 (前端模拟模式)'
    }
  } catch (err: any) {
    testResultSuccess.value = false
    testResultMsg.value = '连接出错: ' + err
  } finally {
    isTesting.value = false
  }
}

async function saveConnection() {
  try {
    if (wails?.SaveConnection) {
      const res = await wails.SaveConnection(formConn.value)
      if (res.success) {
        showConnModal.value = false
        await fetchConnections()
      } else {
        alert(res.message)
      }
    } else {
      showConnModal.value = false
      connections.value.push({ ...formConn.value, id: 'conn_' + Date.now() })
    }
  } catch (e) {
    alert('保存连接失败: ' + e)
  }
}

async function confirmDeleteConnection(conn: ConnectionConfig) {
  if (confirm(`确定要删除连接 "${conn.name}" 吗？`)) {
    if (wails?.DeleteConnection && conn.id) {
      await wails.DeleteConnection(conn.id)
      fetchConnections()
    } else {
      connections.value = connections.value.filter(c => c.id !== conn.id)
    }
  }
}

async function toggleConnection(conn: ConnectionConfig) {
  selectedConn.value = conn
  conn.isExpanded = !conn.isExpanded

  if (conn.isExpanded && (!conn.databases || conn.databases.length === 0)) {
    if (wails?.GetDatabases) {
      const dbs = await wails.GetDatabases(conn)
      conn.databases = dbs || []
    } else {
      conn.databases = ['mysql', 'information_schema', 'performance_schema', 'sys']
    }
  }
}

async function toggleDatabase(conn: ConnectionConfig, dbName: string) {
  selectedConn.value = conn
  selectedDb.value = dbName
  if (!conn.expandedDbs) conn.expandedDbs = {}
  conn.expandedDbs[dbName] = !conn.expandedDbs[dbName]

  if (conn.expandedDbs[dbName] && (!conn.dbTables?.[dbName] || conn.dbTables[dbName].length === 0)) {
    if (!conn.dbTables) conn.dbTables = {}
    if (wails?.GetTables) {
      const tables = await wails.GetTables(conn, dbName)
      conn.dbTables[dbName] = tables || []
    } else {
      conn.dbTables[dbName] = ['users', 'orders', 'products', 'logs']
    }
  }
}

function selectTable(conn: ConnectionConfig, dbName: string, tableName: string) {
  selectedConn.value = conn
  selectedDb.value = dbName
  currentTable.value = tableName
  sqlCode.value = `SELECT * FROM \`${dbName}\`.\`${tableName}\` LIMIT 100;`
  runSQL()
}

async function runSQL() {
  if (!sqlCode.value.trim()) return
  isExecuting.value = true
  currentPage.value = 1

  try {
    if (wails?.ExecuteSQL && selectedConn.value) {
      const res = await wails.ExecuteSQL(selectedConn.value, selectedDb.value, sqlCode.value)
      queryResult.value = res
    } else {
      // Mock execution fallback
      queryResult.value = {
        success: true,
        message: '查询成功，共 2 条记录 (前端预览模式)',
        columns: ['id', 'user_name', 'email', 'created_at'],
        rows: [
          ['1', 'admin', 'admin@linuxnavicat.org', '2026-08-12 10:00:00'],
          ['2', 'developer', 'dev@linuxnavicat.org', '2026-08-12 11:30:00']
        ],
        executionTime: 12
      }
    }
  } catch (e: any) {
    queryResult.value = {
      success: false,
      message: e.toString()
    }
  } finally {
    isExecuting.value = false
  }
}

async function fetchTableDDL() {
  if (!currentTable.value || !selectedConn.value) return
  if (wails?.GetTableDDL) {
    const ddl = await wails.GetTableDDL(selectedConn.value, selectedDb.value, currentTable.value)
    tableDDLContent.value = ddl
  } else {
    tableDDLContent.value = `CREATE TABLE \`${currentTable.value}\` (\n  \`id\` int(11) NOT NULL AUTO_INCREMENT,\n  \`name\` varchar(255) DEFAULT NULL,\n  PRIMARY KEY (\`id\`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
  }
  showDDLModal.value = true
}

async function exportData(format: string) {
  if (!queryResult.value?.columns || !queryResult.value?.rows) return
  let content = ''
  if (wails?.ExportData) {
    content = await wails.ExportData(queryResult.value.columns, queryResult.value.rows, format)
  } else {
    content = JSON.stringify(queryResult.value.rows)
  }

  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `export_${Date.now()}.${format}`
  a.click()
  URL.revokeObjectURL(url)
}
</script>

<style scoped>
.app-container {
  display: flex;
  width: 100vw;
  height: 100vh;
  background: #0f172a;
  color: #f8fafc;
  font-family: system-ui, -apple-system, sans-serif;
  user-select: none;
}

/* Light Theme Overrides */
.app-container.light-theme {
  background: #f8fafc;
  color: #0f172a;
}
.light-theme .sidebar {
  background: #ffffff;
  border-right-color: #e2e8f0;
}
.light-theme .top-header, .light-theme .tree-item:hover, .light-theme .node-header:hover {
  background: #f1f5f9;
}
.light-theme .main-content {
  background: #f8fafc;
}
.light-theme .sql-textarea, .light-theme .grid-wrapper {
  background: #ffffff;
  border-color: #cbd5e1;
  color: #0f172a;
}
.light-theme .grid-table th {
  background: #f1f5f9;
  color: #334155;
}

/* Sidebar Styling */
.sidebar {
  width: 280px;
  background: #1e293b;
  border-right: 1px solid #334155;
  display: flex;
  flex-direction: column;
  padding: 16px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.brand h2 {
  font-size: 1.15rem;
  margin: 0;
  color: #38bdf8;
  font-weight: 700;
}

.version-tag {
  background: #0284c7;
  color: #fff;
  font-size: 0.7rem;
  padding: 2px 6px;
  border-radius: 4px;
}

.btn {
  border: none;
  border-radius: 6px;
  padding: 8px 14px;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: #0284c7;
  color: #ffffff;
  width: 100%;
}
.btn-primary:hover {
  background: #0369a1;
}

.btn-secondary {
  background: #475569;
  color: #fff;
}
.btn-secondary:hover {
  background: #334155;
}

.btn-run {
  background: #16a34a;
  color: #fff;
}
.btn-run:hover {
  background: #15803d;
}

.btn-theme {
  background: #334155;
  color: #f8fafc;
}

.tree-section {
  flex: 1;
  overflow-y: auto;
}

.tree-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.8rem;
  color: #94a3b8;
  margin-bottom: 10px;
}

.btn-icon {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 0.9rem;
}

.tree-list, .tree-sub-list {
  list-style: none;
  padding-left: 0;
  margin: 0;
}
.tree-sub-list {
  padding-left: 14px;
}

.node-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.85rem;
}
.node-header:hover, .node-header.active {
  background: #334155;
}

.tree-arrow {
  font-size: 0.75rem;
  transition: transform 0.2s;
  color: #64748b;
}
.tree-arrow.expanded {
  transform: rotate(90deg);
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}
.status-dot.online {
  background: #22c55e;
}

.node-label {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.node-sub {
  font-size: 0.7rem;
  color: #64748b;
}

.btn-delete {
  background: none;
  border: none;
  color: #64748b;
  cursor: pointer;
  display: none;
}
.node-header:hover .btn-delete {
  display: block;
}
.btn-delete:hover {
  color: #ef4444;
}

.table-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 8px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.8rem;
}
.table-item:hover {
  background: #334155;
}

/* Main Content Area */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.top-header {
  height: 48px;
  background: #1e293b;
  border-bottom: 1px solid #334155;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
}

.tab-list {
  display: flex;
  gap: 4px;
}

.tab-item {
  padding: 6px 16px;
  background: #0f172a;
  border-radius: 6px 6px 0 0;
  font-size: 0.85rem;
  border: 1px solid #334155;
  border-bottom: none;
  color: #38bdf8;
}

.top-controls {
  display: flex;
  gap: 8px;
}

.editor-pane {
  height: 200px;
  padding: 12px;
}

.sql-textarea {
  width: 100%;
  height: 100%;
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 8px;
  color: #38bdf8;
  font-family: 'Fira Code', monospace;
  font-size: 0.95rem;
  padding: 12px;
  box-sizing: border-box;
  resize: none;
}

.status-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px;
  background: #1e293b;
  border-top: 1px solid #334155;
  border-bottom: 1px solid #334155;
  font-size: 0.85rem;
}

.text-success { color: #4ade80; }
.text-danger { color: #f87171; }
.execution-time { color: #94a3b8; margin-left: 8px; }

.export-actions {
  display: flex;
  gap: 6px;
}

.btn-export {
  background: #334155;
  border: none;
  color: #f8fafc;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 0.75rem;
  cursor: pointer;
}
.btn-export:hover {
  background: #475569;
}

.result-pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 12px;
  overflow: hidden;
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #64748b;
  font-size: 0.9rem;
}

.error-box {
  background: #451a1a;
  border: 1px solid #7f1d1d;
  padding: 16px;
  border-radius: 6px;
}
.error-title { font-weight: bold; color: #f87171; }
.error-msg { margin: 8px 0 0 0; color: #fca5a5; font-family: monospace; }

.grid-wrapper {
  flex: 1;
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 8px;
  overflow: auto;
}

.data-grid {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.85rem;
}

.data-grid th {
  background: #334155;
  padding: 8px 12px;
  position: sticky;
  top: 0;
  color: #cbd5e1;
}

.data-grid td {
  padding: 8px 12px;
  border-bottom: 1px solid #334155;
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.col-index { width: 50px; text-align: center; color: #64748b; }

.pagination-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 10px;
  font-size: 0.8rem;
  color: #94a3b8;
}

.page-btn-group {
  display: flex;
  gap: 4px;
}

.btn-page {
  background: #334155;
  border: none;
  color: #fff;
  padding: 4px 8px;
  border-radius: 4px;
  cursor: pointer;
}
.btn-page:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* Modals */
.modal-backdrop {
  position: fixed;
  top: 0; left: 0; width: 100vw; height: 100vh;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-card {
  width: 480px;
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 12px;
  box-shadow: 0 20px 25px -5px rgba(0,0,0,0.5);
  overflow: hidden;
}

.modal-large {
  width: 700px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: #0f172a;
  border-bottom: 1px solid #334155;
}

.modal-header h3 {
  margin: 0; font-size: 1.05rem; color: #38bdf8;
}

.modal-close {
  background: none; border: none; color: #94a3b8; cursor: pointer; font-size: 1.1rem;
}

.modal-body {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 0.8rem; color: #cbd5e1;
}

.form-group input {
  background: #0f172a;
  border: 1px solid #334155;
  padding: 8px 12px;
  border-radius: 6px;
  color: #fff;
  font-size: 0.9rem;
}

.form-row {
  display: flex;
  gap: 12px;
}

.flex-1 { flex: 1; }
.flex-2 { flex: 2; }

.modal-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: #0f172a;
  border-top: 1px solid #334155;
}

.right-actions {
  display: flex;
  gap: 8px;
}

.test-feedback {
  font-size: 0.85rem;
  padding: 8px 12px;
  border-radius: 6px;
  background: #0f172a;
  border: 1px solid #334155;
}

.ddl-code {
  background: #0f172a;
  padding: 16px;
  border-radius: 6px;
  color: #38bdf8;
  font-family: monospace;
  max-height: 400px;
  overflow: auto;
  margin: 0;
}
</style>
