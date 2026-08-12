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

// TestConnection attempts to ping MySQL database
func TestConnection(config ConnectionConfig) (bool, string) {
	if config.Port <= 0 {
		config.Port = 3306
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=5s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Schema,
	)

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
