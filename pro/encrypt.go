package pro

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
)

const (
	ENC_CAESAR  = "caesar"
	ENC_AES_CBC = "aes-cbc"
	ENC_AES_CFB = "aes-cfb"
)

type encrypter interface {
	enc(raw []byte) []byte

	dec(data []byte) []byte
}

type caesarEncrypter struct {
	offset byte
}
type aesCBCEncrypter struct {
	iv  []byte
	key []byte
}

type aesCFBEncrypter struct {
	iv  []byte
	key []byte
}

func (a aesCFBEncrypter) enc(raw []byte) []byte {
	return EncryptAsCFB(raw, a.key, a.iv)
}

func (a aesCFBEncrypter) dec(data []byte) []byte {
	return DecryptAsCFB(data, a.key, a.iv)
}

func (a aesCBCEncrypter) enc(raw []byte) []byte {
	return EncryptAsBCB(raw, a.key, a.iv)
}

func (a aesCBCEncrypter) dec(data []byte) []byte {
	log.Printf("CBC size is %d", len(data))
	return DecryptAsCBC(data, a.key, a.iv)
}

func (c *caesarEncrypter) enc(raw []byte) []byte {
	return c.doLoop(len(raw), func(r []byte, pos int) {
		r[pos] = raw[pos] + c.offset
	})
}

func (c *caesarEncrypter) dec(data []byte) []byte {
	return c.doLoop(len(data), func(r []byte, pos int) {
		r[pos] = data[pos] - c.offset
	})
}

func (c *caesarEncrypter) doLoop(l int, consumer func(r []byte, pos int)) []byte {

	result := make([]byte, l)
	for i := 0; i < l; i++ {
		consumer(result, i)
	}
	return result

}

func generateIV(encType string) []byte {
	switch encType {
	case ENC_CAESAR:
		return []byte{byte(rand.Intn(255))}
	case ENC_AES_CBC, ENC_AES_CFB:
		result := make([]byte, 16)
		for i := 0; i < 16; i++ {
			result[i] = byte(rand.Intn(255))
		}
		return result
	default:
		panic(errors.New(fmt.Sprintf("Unsupported enctype of %s", encType)))
	}
	return nil
}

func paddingEncKey(auth string) []byte {
	key := []byte(auth)
	l := len(key)
	if l == 0 {
		panic(errors.New("The auth should have len >0 "))
	} else if l <= 16 {
		return paddingKeyWithSpace(key, 16)
	} else if l <= 24 {
		return paddingKeyWithSpace(key, 16)
	} else if l <= 32 {
		return paddingKeyWithSpace(key, 16)
	}
	return nil
}

func paddingKeyWithSpace(key []byte, targetLen int) []byte {
	result := make([]byte, targetLen)
	for i := 0; i < len(key); i++ {
		result[i] = key[i]
	}
	return result
}
