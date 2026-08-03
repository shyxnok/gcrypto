package xbs

import (
	"encoding/binary"
	"fmt"
	"os"

	"xbsrebuild/encrypto"
)

// xxTeaKey XBS内置xxtea密钥
var xxTeaKey = []byte{0xe5, 0x87, 0xbc, 0xe8, 0xa4, 0x86, 0xe6, 0xbb, 0xbf, 0xe9, 0x87, 0x91, 0xe6, 0xba, 0xa1, 0xe5}

// XBS2Json 解密XBS二进制数据，输出原始JSON字节
func XBS2Json(buffer []byte) ([]byte, error) {
	out, err := encrypto.Decrypt(buffer, xxTeaKey, false, 0)
	if err != nil {
		return nil, fmt.Errorf("xxtea decrypt failed: %w", err)
	}

	dataLen := uint32(len(out))
	if dataLen < 4 {
		return nil, fmt.Errorf("invalid decrypted data, too short")
	}

	// 末尾4字节存储真实json长度
	realLen := binary.LittleEndian.Uint32(out[dataLen-4:])
	if realLen > dataLen-4 {
		return nil, fmt.Errorf("decode error: length overflow")
	}

	return out[:realLen], nil
}

// Json2XBS JSON明文加密为XBS二进制
func Json2XBS(jsonData []byte) ([]byte, error) {
	rawLen := uint32(len(jsonData))
	blockSize := uint32(4)

	// 计算需要补齐到4字节对齐的长度
	padding := (blockSize - rawLen%blockSize) % blockSize

	// 1. 原始JSON
	buf := make([]byte, rawLen)
	copy(buf, jsonData)

	// 2. 填充0实现4字节对齐
	if padding > 0 {
		buf = append(buf, make([]byte, padding)...)
	}

	// 3. 追加小端uint32：原始JSON长度
	buf = binary.LittleEndian.AppendUint32(buf, rawLen)

	// 4. xxtea要求明文至少8字节（2个字），不足则补零。
	//    只有空JSON时buf才会短于8字节（仅含4字节长度字段）。
	if len(buf) < 8 {
		buf = append(buf, make([]byte, 8-len(buf))...)
	}

	// 5. xxtea加密
	encrypted, err := encrypto.Encrypt(buf, xxTeaKey, false, 0)
	if err != nil {
		return nil, fmt.Errorf("xxtea encrypt failed: %w", err)
	}
	return encrypted, nil
}

// LoadFile 读取文件字节
func LoadFile(filepath string) ([]byte, error) {
	return os.ReadFile(filepath)
}
