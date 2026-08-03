package encrypto

import (
	"bytes"
	"testing"
)

func zeroBytes(n int) []byte {
	return make([]byte, n)
}

var data = []byte("xxtea-go test case")
var key = []byte{82, 253, 252, 7, 33, 130, 101, 79, 22, 63, 95, 15, 154, 98, 29, 114}

var enc = []byte{165, 49, 49, 102, 222, 29, 124, 20, 219, 59, 80, 14, 80, 113, 186, 239, 7, 66, 98, 216}
var hexEnc = "a5313166de1d7c14db3b500e5071baef074262d8"
var b64Enc = "pTExZt4dfBTbO1AOUHG67wdCYtg="

func TestEncrypt(t *testing.T) {
	v, err := Encrypt(data, key, true, 0)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(enc, v) {
		t.Errorf("%+v != %+v", enc, v)
	}
}

func TestDecrypt(t *testing.T) {
	v, err := Decrypt(enc, key, true, 0)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data, v) {
		t.Errorf("%+v != %+v", data, v)
	}
}

func TestEncryptBase64(t *testing.T) {
	v, err := EncryptBase64(data, key, true, 0)
	if err != nil {
		t.Error(err)
	}

	if b64Enc != v {
		t.Errorf("%s != %s", b64Enc, v)
	}
}

func TestDecryptBase64(t *testing.T) {
	v, err := DecryptBase64(b64Enc, key, true, 0)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data, v) {
		t.Errorf("%+v != %+v", data, v)
	}
}

func TestEncryptHex(t *testing.T) {
	v, err := EncryptHex(data, key, true, 0)
	if err != nil {
		t.Error(err)
	}

	if hexEnc != v {
		t.Errorf("%s != %s", hexEnc, v)
	}
}

func TestDecryptHex(t *testing.T) {
	v, err := DecryptHex(hexEnc, key, true, 0)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data, v) {
		t.Errorf("%+v != %+v", data, v)
	}
}

func TestRandom(t *testing.T) {
	for i := 0; i < 2048; i++ {
		data, _ := URandom(i)
		key, _ := URandom(16)

		enc, err := Encrypt(data, key, true, 0)
		if err != nil {
			t.Error(err)
		}

		dec, err := Decrypt(enc, key, true, 0)
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(data, dec) {
			t.Errorf("data=%+v not equal dec=%+v, key=%+v\n", data, dec, key)
		}

		key = zeroBytes(16)
		enc, err = Encrypt(data, key, true, 0)
		if err != nil {
			t.Error(err)
		}
		dec, err = Decrypt(enc, key, true, 0)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(data, dec) {
			t.Errorf("data=%+v not equal dec=%+v, key=%+v\n", data, dec, key)
		}
	}
}

func TestZeroBytes(t *testing.T) {
	for i := 0; i < 2048; i++ {
		data := zeroBytes(i)
		key, _ := URandom(16)

		enc, err := Encrypt(data, key, true, 0)
		if err != nil {
			t.Error(err)
		}
		dec, err := Decrypt(enc, key, true, 0)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(data, dec) {
			t.Errorf("data=%+v not equal dec=%+v, key=%+v\n", data, dec, key)
		}
	}
}

func TestEncryptNoPadding(t *testing.T) {
	key, _ := URandom(16)
	for _, v := range []int{8, 12, 16, 20} {
		data, _ := URandom(v)
		enc, err := Encrypt(data, key, false, 0)
		if err != nil {
			t.Error(err)
		}

		dec, err := Decrypt(enc, key, false, 0)
		if err != nil {
			t.Error(err)
		}

		//t.Logf("log: data=%+v, enc=%+v, dec=%+v", data, enc, dec)
		if !bytes.Equal(data, dec) {
			t.Errorf("data=%+v not equal dec=%+v, key=%+v\n", data, dec, key)
		}
	}
}

func TestEncryptRandomRounds(t *testing.T) {
	key, _ := URandom(16)
	data, _ := URandom(64)
	for i := 1; i < 2048; i++ {
		enc, err := Encrypt(data, key, true, uint32(i))
		if err != nil {
			t.Error(err)
		}

		dec, err := Decrypt(enc, key, true, uint32(i))
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(data, dec) {
			t.Errorf("data=%+v not equal dec=%+v, key=%+v\n", data, dec, key)
		}
	}
}

func TestEncryptNoPaddingZero(t *testing.T) {
	key, _ := URandom(16)
	for _, v := range []int{8, 12, 16, 20} {
		data := zeroBytes(v)

		enc, err := Encrypt(data, key, false, 0)
		if err != nil {
			t.Error(err)
		}

		dec, err := Decrypt(enc, key, false, 0)
		if err != nil {
			t.Error(err)
		}

		//t.Logf("log: data=%+v, enc=%+v, dec=%+v", data, enc, dec)
		if !bytes.Equal(data, dec) {
			t.Errorf("data=%+v not equal dec=%+v, key=%+v\n", data, dec, key)
		}
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	wrongKey := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	out, err := Decrypt(enc, wrongKey, true, 0)
	if err == nil && bytes.Equal(out, data) {
		t.Error("wrong key decrypted the ciphertext to the original plaintext")
	}
}

func TestDecryptCorruptedData(t *testing.T) {
	corrupted := append([]byte(nil), enc...)
	corrupted[0] ^= 0xff
	out, err := Decrypt(corrupted, key, true, 0)
	if err == nil && bytes.Equal(out, data) {
		t.Error("corrupted ciphertext decrypted to the original plaintext")
	}
}

func TestInvalidKeyLength(t *testing.T) {
	shortKey := []byte{1, 2, 3}
	if _, err := Encrypt(data, shortKey, true, 0); err == nil {
		t.Error("expected an error for Encrypt with a non-16-byte key")
	}
	if _, err := Decrypt(enc, shortKey, true, 0); err == nil {
		t.Error("expected an error for Decrypt with a non-16-byte key")
	}
}

func TestEncryptNoPaddingInvalidLength(t *testing.T) {
	for _, n := range []int{0, 2, 3, 4, 6, 7, 9, 11} {
		if _, err := Encrypt(zeroBytes(n), key, false, 0); err == nil {
			t.Errorf("expected an error for nopadding Encrypt of %d bytes", n)
		}
	}
}

func TestDecryptInvalidLength(t *testing.T) {
	for _, n := range []int{0, 4, 5, 6, 9} {
		if _, err := Decrypt(zeroBytes(n), key, true, 0); err == nil {
			t.Errorf("expected an error for Decrypt of %d bytes", n)
		}
	}
}

func TestInvalidBase64(t *testing.T) {
	if _, err := DecryptBase64("!!!not-base64!!!", key, true, 0); err == nil {
		t.Error("expected an error for invalid base64 input")
	}
}

func TestInvalidHex(t *testing.T) {
	if _, err := DecryptHex("zz", key, true, 0); err == nil {
		t.Error("expected an error for invalid hex input")
	}
}
