package xbs

import (
	"bytes"
	"encoding/binary"
	"testing"
	"xbsrebuild/encrypto"
)

// TestRoundTrip 覆盖各种长度的JSON往返加解密，包括此前会失败的nil/空数据。
func TestRoundTrip(t *testing.T) {
	cases := [][]byte{
		nil,
		{},
		[]byte(""),
		[]byte("{}"),
		[]byte("1"),
		[]byte("ab"),
		[]byte("abc"),
		[]byte(`{"a":1}`),
		bytes.Repeat([]byte("x"), 100),
		bytes.Repeat([]byte("x"), 101),
		bytes.Repeat([]byte("x"), 102),
		bytes.Repeat([]byte("x"), 103),
		[]byte(`{"id":123,"name":"张三","tags":["a","b","c"]}`),
	}
	for i, tc := range cases {
		enc, err := Json2XBS(tc)
		if err != nil {
			t.Fatalf("case %d: Json2XBS(%q) error: %v", i, tc, err)
		}
		dec, err := XBS2Json(enc)
		if err != nil {
			t.Fatalf("case %d: XBS2Json error: %v", i, err)
		}
		if !bytes.Equal(dec, tc) {
			t.Errorf("case %d: roundtrip mismatch: got %q want %q", i, dec, tc)
		}
	}
}

// TestEmptyDataCiphertextLength 验证空数据能成功加密，且密文长度合规。
func TestEmptyDataCiphertextLength(t *testing.T) {
	enc, err := Json2XBS(nil)
	if err != nil {
		t.Fatalf("Json2XBS(nil) error: %v", err)
	}
	if len(enc) < 8 {
		t.Errorf("empty data ciphertext too short: %d bytes", len(enc))
	}
	if len(enc)%4 != 0 {
		t.Errorf("ciphertext length not a multiple of 4: %d", len(enc))
	}
}

// TestXBS2JsonGarbage 非法输入（末尾长度字段超界）应报错而不是panic。
func TestXBS2JsonGarbage(t *testing.T) {
	// 8字节垃圾数据，解密后末尾4字节 0x08070605 远超可用长度
	garbage := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	if _, err := XBS2Json(garbage); err == nil {
		t.Error("expected error for garbage input")
	}
}

// TestXBS2JsonLengthOverflow 构造合法密文但长度字段越界，应报错。
func TestXBS2JsonLengthOverflow(t *testing.T) {
	// 明文：4字节内容 + 4字节长度=999（越界）
	buf := []byte{1, 2, 3, 4}
	buf = binary.LittleEndian.AppendUint32(buf, 999)
	enc, err := encrypto.Encrypt(buf, xxTeaKey, false, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := XBS2Json(enc); err == nil {
		t.Error("expected length overflow error")
	}
}
