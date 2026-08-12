package main

import (
	"context"
	"fmt"

	"github.com/your-username/linux-navicat/pkg/db"
)

// App struct
type App struct {
	ctx context.Context
}

type ConnectionResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, Welcome to LinuxNavicat Go!", name)
}

// GetVersion returns current application version
func (a *App) GetVersion() string {
	return "0.2.0-go"
}

// TestConnection tests a MySQL database connection
func (a *App) TestConnection(config db.ConnectionConfig) ConnectionResult {
	success, msg := db.TestConnection(config)
	return ConnectionResult{
		Success: success,
		Message: msg,
	}
}

// SaveConnection saves a connection config to file
func (a *App) SaveConnection(config db.ConnectionConfig) ConnectionResult {
	err := db.GetStorageManager().SaveConnection(config)
	if err != nil {
		return ConnectionResult{
			Success: false,
			Message: fmt.Sprintf("保存配置失败: %v", err),
		}
	}
	return ConnectionResult{
		Success: true,
		Message: "连接保存成功！",
	}
}

// GetConnections retrieves all saved connection configs
func (a *App) GetConnections() []db.ConnectionConfig {
	conns, err := db.GetStorageManager().GetConnections()
	if err != nil {
		return []db.ConnectionConfig{}
	}
	return conns
}

// DeleteConnection removes a connection by ID
func (a *App) DeleteConnection(id string) ConnectionResult {
	err := db.GetStorageManager().DeleteConnection(id)
	if err != nil {
		return ConnectionResult{
			Success: false,
			Message: fmt.Sprintf("删除连接失败: %v", err),
		}
	}
	return ConnectionResult{
		Success: true,
		Message: "连接已删除！",
	}
}

// GetDatabases returns list of databases for a connection
func (a *App) GetDatabases(config db.ConnectionConfig) []string {
	dbs, err := db.GetDatabases(config)
	if err != nil {
		return []string{}
	}
	return dbs
}

// GetTables returns list of tables in a database
func (a *App) GetTables(config db.ConnectionConfig, dbName string) []string {
	tables, err := db.GetTables(config, dbName)
	if err != nil {
		return []string{}
	}
	return tables
}

// GetTableDDL returns DDL statement for a table
func (a *App) GetTableDDL(config db.ConnectionConfig, dbName string, tableName string) string {
	ddl, err := db.GetTableDDL(config, dbName, tableName)
	if err != nil {
		return fmt.Sprintf("获取 DDL 失败: %v", err)
	}
	return ddl
}

// ExecuteSQL executes SQL query and returns formatted result
func (a *App) ExecuteSQL(config db.ConnectionConfig, dbName string, sqlText string) db.QueryResult {
	res, err := db.ExecuteSQL(config, dbName, sqlText)
	if err != nil {
		return db.QueryResult{
			Success: false,
			Message: err.Error(),
		}
	}
	return res
}

// ExportData exports result data as CSV or JSON
func (a *App) ExportData(columns []string, rows [][]string, format string) string {
	if format == "json" {
		return db.ExportToJSON(columns, rows)
	}
	return db.ExportToCSV(columns, rows)
}
