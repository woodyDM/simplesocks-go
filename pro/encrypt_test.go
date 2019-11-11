package pro

import "testing"

func Test_padding_Key(t *testing.T) {
	key := paddingEncKey("haha😄")
	if len(key) != 16 {
		t.Fatal("16")
	}
}

func Test_generate_iv_of_caesar(t *testing.T) {
	iv := generateIV(ENC_CAESAR)
	if len(iv) != 1 {
		t.Fatal("len should = 1")
	}
}

func Test_generate_iv_of_aes_cbc(t *testing.T) {
	generateIvOfAes(t, ENC_AES_CBC)
}

func Test_generate_iv_of_aes_cfb(t *testing.T) {
	generateIvOfAes(t, ENC_AES_CFB)
}

func generateIvOfAes(t *testing.T, encType string) {
	iv := generateIV(encType)
	if len(iv) != 16 {
		t.Fatal("len should = 1")
	}
}
