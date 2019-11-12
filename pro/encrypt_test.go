package pro

import "testing"

func Test_padding_Key_16(t *testing.T) {
	key := paddingEncKey("haha😄")
	if len(key) != 16 {
		t.Fatal("16")
	}
}

func Test_padding_Key_24(t *testing.T) {
	key := paddingEncKey("haha😄1234567890123")
	if len(key) != 24 {
		t.Fatal("24")
	}
	if key[23] != byte(3) {
		t.Fatal("pkcs5 :should be 3")
	}
}

func Test_padding_Key_32(t *testing.T) {
	key := paddingEncKey("haha😄1234567890123456_")
	if len(key) != 32 {
		t.Fatal("32")
	}
	if key[30] != byte(7) {
		t.Fatal("pkcs5 :should be 3")
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
