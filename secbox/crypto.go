package secbox

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

// delta 为 XXTEA 的黄金分割常数，参与每轮加密的扩散。
const delta = 0x9e3779b9

// mx 为 XXTEA 的单步混合函数。
func mx(y, z, p, e, sum uint32, key []uint32) uint32 {
	return ((z>>5 ^ y<<2) + (y>>3 ^ z<<4)) ^ ((sum ^ y) + (key[(p&3)^e] ^ z))
}

// btea 对 uint32 数组执行 XXTEA 变换。n>1 为加密，n<-1 为解密。
func btea(v []uint32, n int, key []uint32, rounds uint32) {
	var i, y, z, p, e, sum uint32
	// Coding Part
	if n > 1 {
		un := uint32(n)
		if rounds == 0 {
			rounds = 6 + 52/un
		}
		z = v[n-1]

		for {
			sum += delta
			e = (sum >> 2) & 3
			for p = 0; p < un-1; p++ {
				y = v[p+1]
				v[p] += mx(y, z, p, e, sum, key)
				z = v[p]
			}

			y = v[0]
			v[n-1] += mx(y, z, p, e, sum, key)
			z = v[n-1]

			i++
			if i > rounds-1 {
				break
			}
		}

	} else if n < -1 { // Decoding Part
		un := uint32(-n)
		if rounds == 0 {
			rounds = 6 + 52/un
		}

		sum = rounds * delta
		y = v[0]

		for {
			e = (sum >> 2) & 3
			for p = un - 1; p > 0; p-- {
				z = v[p-1]
				v[p] -= mx(y, z, p, e, sum, key)
				y = v[p]
			}

			z = v[un-1]
			v[0] -= mx(y, z, p, e, sum, key)
			y = v[0]
			sum -= delta

			i++
			if i > rounds-1 {
				break
			}
		}
	}
}

// bytesToUint32 把字节切片打包为 uint32 数组；padding 为真时按 PKCS#7 填充补齐。
func bytesToUint32(in []byte, inLen int, out []uint32, padding bool) {
	// (i & 3) << 3 -> [0, 8, 16, 24]
	for i := 0; i < inLen; i++ {
		out[i>>2] |= uint32(in[i]) << ((i & 3) << 3)
	}

	if padding {
		pad := 4 - (inLen & 3)
		if inLen < 4 {
			pad = pad + 4
		}

		for i := inLen; i < inLen+pad; i++ {
			out[i>>2] |= uint32(pad) << ((i & 3) << 3)
		}
	}
}

// uint32sToBytes 把 uint32 数组还原为字节切片；padding 为真时校验并去除 PKCS#7 填充。
// 返回去除填充后的有效长度，负数表示填充非法。
func uint32sToBytes(in []uint32, inLen int, out []byte, padding bool) int {
	for i := 0; i < inLen; i++ {
		out[4*i] = byte(in[i] & 0xFF)
		out[4*i+1] = byte((in[i] >> 8) & 0xFF)
		out[4*i+2] = byte((in[i] >> 16) & 0xFF)
		out[4*i+3] = byte((in[i] >> 24) & 0xFF)
	}

	outLen := inLen * 4
	// PKCS#7 unpadding
	if padding {
		pad := int(out[outLen-1])
		outLen -= pad

		if pad < 1 || pad > 8 {
			return -1
		}

		if outLen < 0 {
			return -2
		}

		for i := outLen; i < inLen*4; i++ {
			if int(out[i]) != pad {
				return -3
			}
		}
	}
	return outLen
}

// requireKey 校验 XXTEA 密钥必须为 16 字节。
func requireKey(key []byte) error {
	if len(key) != 16 {
		return errors.New("need a 16-byte key")
	}
	return nil
}

// requireNoPaddingLen 校验无填充模式的数据长度：不小于 8 字节且为 4 的倍数。
func requireNoPaddingLen(dataLen int) error {
	if dataLen < 8 || (dataLen&3) != 0 {
		return errors.New("data length must be a multiple of 4 bytes and must not be less than 8 bytes")
	}
	return nil
}

// URandom 返回指定长度的随机字节。
func URandom(n int) ([]byte, error) {
	if n <= 0 {
		return nil, fmt.Errorf("invalid length n=%d", n)
	}
	token := make([]byte, n)
	_, err := rand.Read(token)
	if err != nil {
		return nil, fmt.Errorf("read random failed: %w", err)
	}
	return token, nil
}

// Encrypt 用密钥加密数据。padding 为真时自动 PKCS#7 填充（支持任意长度）；
// 为假时数据长度必须是不小于 8 字节的 4 的倍数。rounds 为 0 时使用默认轮数。
func Encrypt(data []byte, key []byte, padding bool, rounds uint32) ([]byte, error) {
	var aLen, paddingValue int
	dLen := len(data)

	if err := requireKey(key); err != nil {
		return nil, err
	}

	if padding {
		paddingValue = 1
	} else if err := requireNoPaddingLen(dLen); err != nil {
		return nil, err
	}

	if dLen < 4 {
		aLen = 2
	} else {
		aLen = dLen>>2 + paddingValue
	}

	d := make([]uint32, aLen)
	k := make([]uint32, 4)
	bytesToUint32(data, dLen, d, padding)
	bytesToUint32(key, 16, k, false)
	btea(d, aLen, k, rounds)

	retBuf := make([]byte, aLen<<2)
	_ = uint32sToBytes(d, aLen, retBuf, false)
	return retBuf, nil
}

// Decrypt 用密钥解密数据，参数语义与 Encrypt 一致。
func Decrypt(data []byte, key []byte, padding bool, rounds uint32) ([]byte, error) {
	var rc int
	dLen := len(data)

	if err := requireKey(key); err != nil {
		return nil, err
	}

	if (dLen&3) != 0 || dLen < 8 {
		return nil, errors.New("invalid data, data length is not a multiple of 4, or less than 8")
	}

	if !padding {
		if err := requireNoPaddingLen(dLen); err != nil {
			return nil, err
		}
	}

	aLen := dLen / 4
	d := make([]uint32, aLen)
	k := make([]uint32, 4)
	bytesToUint32(data, dLen, d, false)
	bytesToUint32(key, 16, k, false)
	btea(d, -aLen, k, rounds)

	refBuf := make([]byte, dLen)
	rc = uint32sToBytes(d, aLen, refBuf, padding)

	if padding {
		if rc >= 0 {
			refBuf = refBuf[:rc]
		} else {
			return nil, errors.New("invalid data, illegal PKCS#7 padding. Could be using a wrong key")
		}
	}

	return refBuf, nil
}

// EncryptBase64 加密后返回 base64 编码结果。
func EncryptBase64(data []byte, key []byte, padding bool, rounds uint32) (string, error) {
	v, err := Encrypt(data, key, padding, rounds)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(v), nil
}

// DecryptBase64 解密 base64 编码的数据。
func DecryptBase64(b64Str string, key []byte, padding bool, rounds uint32) ([]byte, error) {
	dataBytes, err := base64.StdEncoding.DecodeString(b64Str)
	if err != nil {
		return nil, err
	}

	v, err := Decrypt(dataBytes, key, padding, rounds)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// EncryptHex 加密后返回十六进制编码结果。
func EncryptHex(data []byte, key []byte, padding bool, rounds uint32) (string, error) {
	v, err := Encrypt(data, key, padding, rounds)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(v), nil
}

// DecryptHex 解密十六进制编码的数据。
func DecryptHex(hexStr string, key []byte, padding bool, rounds uint32) ([]byte, error) {
	dataBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}

	v, err := Decrypt(dataBytes, key, padding, rounds)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// textKey 为 EncryptText/DecryptText 使用的固定密钥，包初始化时计算一次。
var textKey = func() []byte {
	key := []byte{
		0xa4, 0x59, 0x94, 0xad,
		0x88, 0x31, 0xb2, 0x1a,
		0xd7, 0xc4, 0x43, 0x8f,
		0x8f, 0x51, 0xe3, 0x29,
	}
	for i := range key {
		key[i] ^= 0x47
	}
	return key
}()

// EncryptText 通用文本/文件加密，支持任意长度二进制，自动 PKCS7 填充，普通人打开为乱码。
func EncryptText(data []byte) ([]byte, error) {
	return Encrypt(data, textKey, true, 0)
}

// DecryptText 通用文本/文件解密。
func DecryptText(data []byte) ([]byte, error) {
	return Decrypt(data, textKey, true, 0)
}
