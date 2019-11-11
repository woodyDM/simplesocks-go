package pro

import (
	"bytes"
	"testing"
)

func Test_aes_cbc(t *testing.T) {
	plain := "haha哈哈😄jofjoeifjqowejfo2ij3o4j1o23j4o12j4o12j4o12j34"
	key := bytes.Repeat([]byte{1, 2, 3, 4}, 8)
	iv := bytes.Repeat([]byte{1, 2, 5, 6}, 4)

	encrypted := EncryptAsBCB([]byte(plain), key, iv)
	decrypted := DecryptAsCBC(encrypted, key, iv)
	dec := string(decrypted)
	if dec != plain {
		t.Fatal("fail.")
	}
}

func Test_aes_cfb(t *testing.T) {

	plain := "haha哈哈😄jofjoeifjqowejfo2ij3o4j1o23j4o12j4o12j4o12j34"
	key := bytes.Repeat([]byte{1, 2, 3, 4}, 6)
	iv := bytes.Repeat([]byte{1, 2, 5, 6}, 4)

	encrypted := EncryptAsCFB([]byte(plain), key, iv)
	decrypted := DecryptAsCFB(encrypted, key, iv)
	dec := string(decrypted)
	if dec != plain {
		t.Fatal("fail.")
	}
}
