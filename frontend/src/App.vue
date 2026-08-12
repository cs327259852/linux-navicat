<template>
  <div class="app-container">
    <!-- Left Sidebar -->
    <aside class="sidebar">
      <div class="brand">
        <span class="icon">🐬</span>
        <h2>LinuxNavicat</h2>
        <span class="badge">Go + Wails</span>
      </div>

      <div class="action-bar">
        <button class="btn btn-primary" @click="showNewConnDialog = true">+ 新建连接</button>
      </div>

      <div class="tree-section">
        <div class="tree-header">数据库连接 (Connections)</div>
        <ul class="tree-list">
          <li v-for="conn in connections" :key="conn.id" class="tree-item" :class="{ active: activeConn === conn.id }" @click="selectConnection(conn.id)">
            <span class="status-dot online"></span>
            <span class="conn-name">{{ conn.name }}</span>
            <span class="conn-host">{{ conn.host }}:{{ conn.port }}</span>
          </li>
        </ul>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="main-content">
      <header class="top-header">
        <div class="tabs">
          <div class="tab active">SQL 查询标签 (Query 1)</div>
        </div>
        <div class="header-actions">
          <button class="btn btn-run" @click="executeSQL">▶ 执行 (Execute)</button>
        </div>
      </header>

      <div class="editor-pane">
        <textarea v-model="sqlCode" class="sql-input" placeholder="输入 SQL 语句 (例: SELECT * FROM users;)"></textarea>
      </div>

      <div class="result-pane">
        <div class="pane-title">执行结果 (Query Results)</div>
        <div class="grid-table-container">
          <table class="grid-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Host</th>
                <th>Database</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>1</td>
                <td>127.0.0.1:3306</td>
                <td>local_db</td>
                <td><span class="tag tag-success">Connected</span></td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const activeConn = ref('1')
const sqlCode = ref('SELECT * FROM information_schema.tables LIMIT 10;')
const showNewConnDialog = ref(false)

const connections = ref([
  { id: '1', name: 'Local MySQL', host: '127.0.0.1', port: 3306, db: 'mysql' }
])

function selectConnection(id: string) {
  activeConn.value = id
}

function executeSQL() {
  alert('执行 SQL: ' + sqlCode.value)
}
</script>

<style scoped>
.app-container {
  display: flex;
  width: 100vw;
  height: 100vh;
  background: #0f172a;
  color: #f8fafc;
}

.sidebar {
  width: 260px;
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
  margin-bottom: 20px;
}

.brand h2 {
  font-size: 1.1rem;
  margin: 0;
  color: #38bdf8;
}

.badge {
  background: #0284c7;
  font-size: 0.7rem;
  padding: 2px 6px;
  border-radius: 4px;
}

.action-bar {
  margin-bottom: 16px;
}

.btn {
  border: none;
  border-radius: 6px;
  padding: 8px 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-primary {
  background: #2563eb;
  color: #fff;
  width: 100%;
}
.btn-primary:hover {
  background: #1d4ed8;
}

.btn-run {
  background: #16a34a;
  color: #fff;
}
.btn-run:hover {
  background: #15803d;
}

.tree-section {
  flex: 1;
}

.tree-header {
  font-size: 0.8rem;
  color: #94a3b8;
  text-transform: uppercase;
  margin-bottom: 10px;
}

.tree-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.tree-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.2s;
}

.tree-item:hover, .tree-item.active {
  background: #334155;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-dot.online {
  background: #22c55e;
}

.conn-name {
  font-size: 0.9rem;
  flex: 1;
}

.conn-host {
  font-size: 0.75rem;
  color: #64748b;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #0f172a;
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

.tabs {
  display: flex;
  gap: 4px;
}

.tab {
  padding: 6px 12px;
  background: #0f172a;
  border-radius: 6px 6px 0 0;
  font-size: 0.85rem;
  border: 1px solid #334155;
  border-bottom: none;
}

.editor-pane {
  height: 220px;
  padding: 12px;
  background: #0f172a;
  border-bottom: 1px solid #334155;
}

.sql-input {
  width: 100%;
  height: 100%;
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 6px;
  color: #f8fafc;
  font-family: monospace;
  font-size: 0.95rem;
  padding: 12px;
  box-sizing: border-box;
  resize: none;
}
.sql-input:focus {
  outline: 1px solid #38bdf8;
}

.result-pane {
  flex: 1;
  padding: 12px;
  display: flex;
  flex-direction: column;
}

.pane-title {
  font-size: 0.85rem;
  color: #94a3b8;
  margin-bottom: 8px;
}

.grid-table-container {
  flex: 1;
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 6px;
  overflow: auto;
}

.grid-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.85rem;
}

.grid-table th {
  background: #334155;
  padding: 8px 12px;
  color: #cbd5e1;
}

.grid-table td {
  padding: 8px 12px;
  border-bottom: 1px solid #334155;
}

.tag {
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.75rem;
}

.tag-success {
  background: #14532d;
  color: #4ade80;
}
</style>
