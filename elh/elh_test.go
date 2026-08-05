package elh

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEncryptFileSuffixHeaderAndName(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "report.docx")
	content := []byte("hello world 你好")
	if err := os.WriteFile(src, content, 0644); err != nil {
		t.Fatal(err)
	}

	dst, err := EncryptFile(src)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dir, "report.lh")
	if dst != want {
		t.Fatalf("dst = %q, want %q", dst, want)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	// 文件头：1 字节后缀长度 + 后缀 + 密文
	if int(data[0]) != len(".docx") || string(data[1:1+data[0]]) != ".docx" {
		t.Fatalf("bad header: len=%d ext=%q", data[0], data[1:1+data[0]])
	}
	if len(data) <= 1+int(data[0]) {
		t.Fatalf("no encrypted payload")
	}
	if strings.Contains(string(data), string(content)) {
		t.Fatalf("plaintext leaked into .lh file")
	}
}

func TestEncryptFileNoExtension(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "notes")
	if err := os.WriteFile(src, []byte("abc"), 0644); err != nil {
		t.Fatal(err)
	}
	dst, err := EncryptFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(dst) != "notes.lh" {
		t.Fatalf("dst = %q, want notes.lh", dst)
	}
}

func TestEncryptFileOverwritesExisting(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.hxl")
	if err := os.WriteFile(src, []byte("new content"), 0644); err != nil {
		t.Fatal(err)
	}
	old := filepath.Join(dir, "a.lh")
	if err := os.WriteFile(old, []byte("old garbage"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := EncryptFile(src); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(old)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) == "old garbage" {
		t.Fatalf("existing a.lh was not replaced")
	}
}

func TestEncryptFileRejectsLHSuffix(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "already.lh")
	if err := os.WriteFile(src, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := EncryptFile(src); err == nil {
		t.Fatalf("expected error for .lh input")
	}
}

func TestEncryptFileExtCustomSuffix(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "data.hxl")
	if err := os.WriteFile(src, []byte("custom suffix"), 0644); err != nil {
		t.Fatal(err)
	}

	// ".enc" 与 "enc" 都应可用
	for _, ext := range []string{"enc", ".enc"} {
		dst, err := EncryptFileExt(src, ext)
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(dir, "data.enc")
		if dst != want {
			t.Fatalf("ext %q: dst = %q, want %q", ext, dst, want)
		}
		data, err := os.ReadFile(dst)
		if err != nil {
			t.Fatal(err)
		}
		if int(data[0]) != len(".hxl") || string(data[1:1+data[0]]) != ".hxl" {
			t.Fatalf("ext %q: bad header", ext)
		}
	}
}

func TestEncryptFileExtRejectsBadSuffix(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "data.hxl")
	if err := os.WriteFile(src, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	for _, ext := range []string{"", ".", "a/b", "a\\b"} {
		if _, err := EncryptFileExt(src, ext); err == nil {
			t.Fatalf("expected error for suffix %q", ext)
		}
	}
}
