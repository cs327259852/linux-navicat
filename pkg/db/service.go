package db

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectionConfig represents a database connection profile
type ConnectionConfig struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Schema    string `json:"schema"`
	CreatedAt string `json:"createdAt,omitempty"`
}

// QueryResult represents the outcome of a SQL execution
type QueryResult struct {
	Success        bool       `json:"success"`
	Message        string     `json:"message"`
	Columns        []string   `json:"columns"`
	Rows           [][]string `json:"rows"`
	AffectedRows   int64      `json:"affectedRows"`
	ExecutionTime  int64      `json:"executionTime"` // in milliseconds
}

type StorageManager struct {
	filePath string
	mu       sync.RWMutex
	secret   []byte
}

var globalStorage *StorageManager
var once sync.Once

// GetStorageManager returns single instance of StorageManager
func GetStorageManager() *StorageManager {
	once.Do(func() {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "."
		}
		dir := filepath.Join(homeDir, ".linuxnavicat")
		_ = os.MkdirAll(dir, 0755)

		globalStorage = &StorageManager{
			filePath: filepath.Join(dir, "connections.json"),
			secret:   []byte("LinuxNavicatSafeKey2026Secret!!!"), // 32 bytes AES key
		}
	})
	return globalStorage
}

func getDSN(config ConnectionConfig, dbName string) string {
	if config.Port <= 0 {
		config.Port = 3306
	}
	schema := config.Schema
	if dbName != "" {
		schema = dbName
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=5s&parseTime=true",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		schema,
	)
}

// TestConnection attempts to ping MySQL database
func TestConnection(config ConnectionConfig) (bool, string) {
	dsn := getDSN(config, "")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return false, fmt.Sprintf("创建连接失败: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return false, fmt.Sprintf("连接失败: %v", err)
	}

	return true, "数据库连接成功！"
}

// GetDatabases returns list of available databases
func GetDatabases(config ConnectionConfig) ([]string, error) {
	dsn := getDSN(config, "")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err == nil {
			databases = append(databases, dbName)
		}
	}
	return databases, nil
}

// GetTables returns list of tables in specified database
func GetTables(config ConnectionConfig, dbName string) ([]string, error) {
	dsn := getDSN(config, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err == nil {
			tables = append(tables, tableName)
		}
	}
	return tables, nil
}

// GetTableDDL returns DDL for a given table
func GetTableDDL(config ConnectionConfig, dbName string, tableName string) (string, error) {
	dsn := getDSN(config, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return "", err
	}
	defer db.Close()

	query := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)
	row := db.QueryRow(query)

	var name, ddl string
	if err := row.Scan(&name, &ddl); err != nil {
		return "", err
	}
	return ddl, nil
}

// ExecuteSQL executes arbitrary SQL statement
func ExecuteSQL(config ConnectionConfig, dbName string, sqlText string) (QueryResult, error) {
	startTime := time.Now()
	dsn := getDSN(config, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return QueryResult{Success: false, Message: fmt.Sprintf("连接失败: %v", err)}, nil
	}
	defer db.Close()

	sqlTextTrimmed := strings.TrimSpace(sqlText)
	isSelect := strings.HasPrefix(strings.ToUpper(sqlTextTrimmed), "SELECT") ||
		strings.HasPrefix(strings.ToUpper(sqlTextTrimmed), "SHOW") ||
		strings.HasPrefix(strings.ToUpper(sqlTextTrimmed), "DESCRIBE") ||
		strings.HasPrefix(strings.ToUpper(sqlTextTrimmed), "EXPLAIN")

	if !isSelect {
		res, err := db.Exec(sqlText)
		elapsed := time.Since(startTime).Milliseconds()
		if err != nil {
			return QueryResult{
				Success:       false,
				Message:       fmt.Sprintf("执行错误: %v", err),
				ExecutionTime: elapsed,
			}, nil
		}
		affected, _ := res.RowsAffected()
		return QueryResult{
			Success:       true,
			Message:       fmt.Sprintf("执行成功，受影响行数: %d", affected),
			AffectedRows:  affected,
			ExecutionTime: elapsed,
		}, nil
	}

	rows, err := db.Query(sqlText)
	elapsed := time.Since(startTime).Milliseconds()
	if err != nil {
		return QueryResult{
			Success:       false,
			Message:       fmt.Sprintf("查询失败: %v", err),
			ExecutionTime: elapsed,
		}, nil
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return QueryResult{Success: false, Message: err.Error(), ExecutionTime: elapsed}, nil
	}

	var resultRows [][]string
	rawValues := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for i := range rawValues {
		valuePtrs[i] = &rawValues[i]
	}

	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			continue
		}
		rowStr := make([]string, len(cols))
		for i, val := range rawValues {
			if val == nil {
				rowStr[i] = "NULL"
			} else if b, ok := val.([]byte); ok {
				rowStr[i] = string(b)
			} else {
				rowStr[i] = fmt.Sprintf("%v", val)
			}
		}
		resultRows = append(resultRows, rowStr)
	}

	return QueryResult{
		Success:       true,
		Message:       fmt.Sprintf("查询成功，共 %d 条记录", len(resultRows)),
		Columns:       cols,
		Rows:          resultRows,
		AffectedRows:  int64(len(resultRows)),
		ExecutionTime: elapsed,
	}, nil
}

// ExportToCSV converts query results to CSV format
func ExportToCSV(columns []string, rows [][]string) string {
	var builder strings.Builder
	builder.WriteString(strings.Join(columns, ",") + "\n")
	for _, row := range rows {
		escapedRow := make([]string, len(row))
		for i, val := range row {
			escapedRow[i] = fmt.Sprintf("\"%s\"", strings.ReplaceAll(val, "\"", "\"\""))
		}
		builder.WriteString(strings.Join(escapedRow, ",") + "\n")
	}
	return builder.String()
}

// ExportToJSON converts query results to JSON array string
func ExportToJSON(columns []string, rows [][]string) string {
	var list []map[string]string
	for _, row := range rows {
		item := make(map[string]string)
		for i, col := range columns {
			if i < len(row) {
				item[col] = row[i]
			}
		}
		list = append(list, item)
	}
	data, _ := json.MarshalIndent(list, "", "  ")
	return string(data)
}

// SaveConnection saves or updates a connection config to file
func (sm *StorageManager) SaveConnection(config ConnectionConfig) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	connections, err := sm.loadUnlocked()
	if err != nil {
		connections = []ConnectionConfig{}
	}

	// Encrypt password before saving
	if config.Password != "" {
		encryptedPwd, err := sm.encrypt(config.Password)
		if err == nil {
			config.Password = encryptedPwd
		}
	}

	if config.ID == "" {
		config.ID = fmt.Sprintf("conn_%d", time.Now().UnixNano())
		config.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
		connections = append(connections, config)
	} else {
		found := false
		for i, c := range connections {
			if c.ID == config.ID {
				connections[i] = config
				found = true
				break
			}
		}
		if !found {
			connections = append(connections, config)
		}
	}

	data, err := json.MarshalIndent(connections, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(sm.filePath, data, 0600)
}

// GetConnections loads all connection configs with decrypted passwords
func (sm *StorageManager) GetConnections() ([]ConnectionConfig, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	connections, err := sm.loadUnlocked()
	if err != nil {
		return []ConnectionConfig{}, nil
	}

	// Decrypt passwords
	for i := range connections {
		if connections[i].Password != "" {
			decrypted, err := sm.decrypt(connections[i].Password)
			if err == nil {
				connections[i].Password = decrypted
			}
		}
	}

	return connections, nil
}

// DeleteConnection removes a connection by ID
func (sm *StorageManager) DeleteConnection(id string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	connections, err := sm.loadUnlocked()
	if err != nil {
		return nil
	}

	var updated []ConnectionConfig
	for _, c := range connections {
		if c.ID != id {
			updated = append(updated, c)
		}
	}

	data, err := json.MarshalIndent(updated, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(sm.filePath, data, 0600)
}

func (sm *StorageManager) loadUnlocked() ([]ConnectionConfig, error) {
	if _, err := os.Stat(sm.filePath); os.IsNotExist(err) {
		return []ConnectionConfig{}, nil
	}

	data, err := os.ReadFile(sm.filePath)
	if err != nil {
		return nil, err
	}

	var connections []ConnectionConfig
	if err := json.Unmarshal(data, &connections); err != nil {
		return nil, err
	}

	return connections, nil
}

func (sm *StorageManager) encrypt(text string) (string, error) {
	block, err := aes.NewCipher(sm.secret)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(text))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (sm *StorageManager) decrypt(cryptoText string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(sm.secret)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
