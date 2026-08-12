package db

import (
	"testing"
)

func TestEncryptionDecryption(t *testing.T) {
	sm := GetStorageManager()
	rawPwd := "SecretPassword123!@#"

	encrypted, err := sm.encrypt(rawPwd)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	if encrypted == rawPwd {
		t.Fatalf("加密输出不能与原文相同")
	}

	decrypted, err := sm.decrypt(encrypted)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}

	if decrypted != rawPwd {
		t.Fatalf("期望解密文本 '%s'，但得到 '%s'", rawPwd, decrypted)
	}
}
