package dlh

import (
	"os"
	"path/filepath"
	"testing"

	"gcrypto/elh"
)

func TestRoundTripWithExtension(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "photo.hxl")
	content := []byte("binary \x00\x01\x02 payload 机密")
	if err := os.WriteFile(src, content, 0644); err != nil {
		t.Fatal(err)
	}

	enc, err := elh.EncryptFile(src)
	if err != nil {
		t.Fatal(err)
	}
	dec, err := DecryptFile(enc)
	if err != nil {
		t.Fatal(err)
	}

	want := filepath.Join(dir, "photo.hxl")
	if dec != want {
		t.Fatalf("dec = %q, want %q", dec, want)
	}
	got, err := os.ReadFile(dec)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Fatalf("round trip mismatch")
	}
}

func TestRoundTripNoExtension(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "notes")
	content := []byte("plain file without extension")
	if err := os.WriteFile(src, content, 0644); err != nil {
		t.Fatal(err)
	}

	enc, err := elh.EncryptFile(src)
	if err != nil {
		t.Fatal(err)
	}
	dec, err := DecryptFile(enc)
	if err != nil {
		t.Fatal(err)
	}
	if dec != filepath.Join(dir, "notes") {
		t.Fatalf("dec = %q, want %q", dec, filepath.Join(dir, "notes"))
	}
	got, err := os.ReadFile(dec)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Fatalf("round trip mismatch")
	}
}

func TestDecryptReplacesExistingFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "doc.hxl")
	if err := os.WriteFile(src, []byte("real"), 0644); err != nil {
		t.Fatal(err)
	}
	enc, err := elh.EncryptFile(src)
	if err != nil {
		t.Fatal(err)
	}

	// 还原前原文件已存在（内容是旧版本），解密时应删除旧的再写回
	if err := os.WriteFile(src, []byte("stale"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := DecryptFile(enc); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "real" {
		t.Fatalf("existing file was not replaced: %q", got)
	}
}

func TestDecryptAnySuffix(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "data.hxl")
	if err := os.WriteFile(src, []byte("any suffix works"), 0644); err != nil {
		t.Fatal(err)
	}

	enc, err := elh.EncryptFileExt(src, "enc")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(enc) != "data.enc" {
		t.Fatalf("enc = %q, want data.enc", enc)
	}
	dec, err := DecryptFile(enc)
	if err != nil {
		t.Fatal(err)
	}
	if dec != filepath.Join(dir, "data.hxl") {
		t.Fatalf("dec = %q, want %q", dec, filepath.Join(dir, "data.hxl"))
	}
	got, err := os.ReadFile(dec)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "any suffix works" {
		t.Fatalf("round trip mismatch")
	}
}

func TestDecryptRejectsBadHeader(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "bad.lh")
	// 头部声称后缀长 200 字节，但文件只有 5 字节
	if err := os.WriteFile(src, []byte{200, 'a', 'b', 'c', 'd'}, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := DecryptFile(src); err == nil {
		t.Fatalf("expected error for bad header")
	}
}
