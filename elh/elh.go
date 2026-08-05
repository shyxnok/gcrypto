// Package elh 提供 .lh 加密文件的生成：原文件后缀写入文件头，输出同名 .lh 文件。
package elh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gcrypto/secbox"
)

// maxExtLen 限制文件头中后缀的最大长度，防止解析恶意构造的文件。
const maxExtLen = 64

// EncryptFile 加密 srcPath 为 .lh 文件，等价于 EncryptFileExt(srcPath, "lh")。
func EncryptFile(srcPath string) (string, error) {
	return EncryptFileExt(srcPath, "lh")
}

// EncryptFileExt 加密 srcPath 为指定后缀的文件（ext 可为 "lh" 或 ".lh"）：
//   - 把原文件的后缀（如 ".hxl"）以 [1 字节长度][后缀] 写入文件头
//   - 原文件名后缀替换为 ext，如 my.hxl + "lh" -> my.lh
//   - 目标已存在时直接删除旧的，再写入新文件
//
// 返回生成的文件路径。
func EncryptFileExt(srcPath, ext string) (string, error) {
	raw, err := os.ReadFile(srcPath)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", srcPath, err)
	}

	origExt := filepath.Ext(srcPath)
	if origExt == ".lh" {
		return "", fmt.Errorf("%s already ends with .lh", srcPath)
	}
	if len(origExt) > maxExtLen {
		return "", fmt.Errorf("extension too long: %d bytes", len(origExt))
	}

	ext = strings.TrimPrefix(strings.TrimSpace(ext), ".")
	if ext == "" || strings.ContainsAny(ext, `/\`) {
		return "", fmt.Errorf("invalid suffix %q", ext)
	}

	encrypted, err := secbox.EncryptText(raw)
	if err != nil {
		return "", fmt.Errorf("encrypt %s: %w", srcPath, err)
	}

	payload := make([]byte, 0, 1+len(origExt)+len(encrypted))
	payload = append(payload, byte(len(origExt)))
	payload = append(payload, origExt...)
	payload = append(payload, encrypted...)

	base := strings.TrimSuffix(filepath.Base(srcPath), origExt)
	dst := filepath.Join(filepath.Dir(srcPath), base+"."+ext)
	if err := overwrite(dst, payload); err != nil {
		return "", fmt.Errorf("write %s: %w", dst, err)
	}
	return dst, nil
}

// overwrite 删除已存在的同名文件（忽略“不存在”错误）后写入数据。
func overwrite(dst string, data []byte) error {
	if err := os.Remove(dst); err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
