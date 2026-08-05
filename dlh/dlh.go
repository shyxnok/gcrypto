// Package dlh 提供 .lh 加密文件的还原：读取文件头中的原后缀，解密后恢复原文件名。
package dlh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gcrypto/secbox"
)

// maxExtLen 必须与 elh 包保持一致：限制文件头后缀最大长度。
const maxExtLen = 64

// DecryptFile 解密 srcPath（任意后缀的加密文件）：
//   - 从文件头读出原后缀（如 ".hxl"）
//   - 输出文件名恢复为去掉当前后缀后加上原后缀，如 my.enc -> my.hxl
//   - 目标已存在时直接删除旧的，再写入还原文件
//
// 返回还原的文件路径。
func DecryptFile(srcPath string) (string, error) {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", srcPath, err)
	}
	if len(data) < 2 {
		return "", fmt.Errorf("invalid .lh file: too short")
	}

	extLen := int(data[0])
	if extLen > maxExtLen || 1+extLen >= len(data) {
		return "", fmt.Errorf("invalid .lh file: bad extension header")
	}
	ext := string(data[1 : 1+extLen])

	raw, err := secbox.DecryptText(data[1+extLen:])
	if err != nil {
		return "", fmt.Errorf("decrypt %s: %w", srcPath, err)
	}

	base := strings.TrimSuffix(filepath.Base(srcPath), filepath.Ext(srcPath))
	dst := filepath.Join(filepath.Dir(srcPath), base+ext)
	if dst == srcPath {
		return "", fmt.Errorf("refusing to overwrite source file %s", srcPath)
	}
	if err := overwrite(dst, raw); err != nil {
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
