package secbox

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestRoundTrip 覆盖各种长度的文件加解密往返，验证 EncryptFile + DecryptFile 后内容与源文件一致。
func TestRoundTrip(t *testing.T) {
	cases := [][]byte{
		[]byte(""),
		[]byte("a"),
		[]byte("ab"),
		[]byte("abc"),
		[]byte("hello world"),
		[]byte("张三，中文内容测试"),
		{0x00, 0x01, 0x02, 0xFF},
		bytes.Repeat([]byte("x"), 100),
		bytes.Repeat([]byte("x"), 101),
		bytes.Repeat([]byte("x"), 102),
		bytes.Repeat([]byte("x"), 103),
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("case%d_len%d", i, len(tc)), func(t *testing.T) {
			dir := t.TempDir()
			src := filepath.Join(dir, "src.bin")
			enc := filepath.Join(dir, "enc.bin")
			dst := filepath.Join(dir, "dst.bin")

			if err := os.WriteFile(src, tc, 0644); err != nil {
				t.Fatalf("write src: %v", err)
			}

			if err := EncryptFile(src, enc); err != nil {
				t.Fatalf("EncryptFile(%q, %q) error: %v", src, enc, err)
			}

			// 密文必须非空，且不能与明文相同。
			cipher, err := os.ReadFile(enc)
			if err != nil {
				t.Fatalf("read cipher: %v", err)
			}
			if len(cipher) == 0 {
				t.Error("ciphertext is empty")
			}
			if bytes.Equal(cipher, tc) {
				t.Error("ciphertext must not equal plaintext")
			}

			if err := DecryptFile(enc, dst); err != nil {
				t.Fatalf("DecryptFile(%q, %q) error: %v", enc, dst, err)
			}

			got, err := os.ReadFile(dst)
			if err != nil {
				t.Fatalf("read dst: %v", err)
			}
			if !bytes.Equal(got, tc) {
				t.Errorf("roundtrip mismatch: got %q want %q", got, tc)
			}
		})
	}
}

// TestEncryptFileSourceMissing 源文件不存在应返回错误。
func TestEncryptFileSourceMissing(t *testing.T) {
	dir := t.TempDir()
	err := EncryptFile(filepath.Join(dir, "nope.txt"), filepath.Join(dir, "out.bin"))
	if err == nil {
		t.Error("expected error for missing source file")
	}
}

// TestDecryptFileSourceMissing 源文件不存在应返回错误。
func TestDecryptFileSourceMissing(t *testing.T) {
	dir := t.TempDir()
	err := DecryptFile(filepath.Join(dir, "nope.bin"), filepath.Join(dir, "out.txt"))
	if err == nil {
		t.Error("expected error for missing source file")
	}
}

// TestEncryptFileBadDest 目标目录不存在时写文件失败，应返回错误。
func TestEncryptFileBadDest(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.bin")
	if err := os.WriteFile(src, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	err := EncryptFile(src, filepath.Join(dir, "no", "such", "dir", "out.bin"))
	if err == nil {
		t.Error("expected error for missing destination directory")
	}
}

// TestDecryptFileGarbage 非法长度数据（非4字节倍数）解密应报错而不是 panic。
func TestDecryptFileGarbage(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "bad.bin")
	dst := filepath.Join(dir, "out.txt")
	if err := os.WriteFile(src, []byte{1, 2, 3, 4, 5, 6, 7}, 0644); err != nil {
		t.Fatal(err)
	}
	if err := DecryptFile(src, dst); err == nil {
		t.Error("expected error for invalid length garbage")
	}
}

// TestDecryptFileEncryptedByEncryptFile 手工构造一份明文加密后再解密，确保数据可还原且目标文件已生成。
func TestDecryptFileEncryptedByEncryptFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	enc := filepath.Join(dir, "enc.bin")
	dst := filepath.Join(dir, "dec.txt")

	payload := []byte("可逆对称加密测试 payload-42")
	if err := os.WriteFile(src, payload, 0644); err != nil {
		t.Fatal(err)
	}
	if err := EncryptFile(src, enc); err != nil {
		t.Fatalf("EncryptFile error: %v", err)
	}
	if err := DecryptFile(enc, dst); err != nil {
		t.Fatalf("DecryptFile error: %v", err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read dec: %v", err)
	}
	if !bytes.Equal(got, payload) {
		t.Errorf("decrypted mismatch: got %q want %q", got, payload)
	}
}
