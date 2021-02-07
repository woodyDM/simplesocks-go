package pkg

import (
	"bytes"
	"errors"
	"fmt"
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

func (a aesCBCEncrypter) enc(raw []byte) []byte {
	return EncryptAsBCB(raw, a.key, a.iv)
}

func (a aesCBCEncrypter) dec(data []byte) []byte {
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

func generateIV(encType string) ([]byte, error) {
	switch encType {
	case ENC_CAESAR:
		return []byte{byte(rand.Intn(255))}, nil
	case ENC_AES_CBC:
		result := make([]byte, 16)
		for i := 0; i < 16; i++ {
			result[i] = byte(rand.Intn(255))
		}
		return result, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported enctype of %s", encType))
	}
}

func paddingEncKey(auth string) []byte {
	key := []byte(auth)
	l := len(key)
	if l == 0 {
		panic(errors.New("The auth should have len >0 "))
	} else if l <= 16 {
		return paddingKeyUsingPkcs5(key, 16)
	} else if l <= 24 {
		return paddingKeyUsingPkcs5(key, 24)
	} else if l <= 32 {
		return paddingKeyUsingPkcs5(key, 32)
	}
	return nil
}

func paddingKeyUsingPkcs5(key []byte, targetLen int) []byte {
	left := targetLen - len(key)
	result := make([]byte, targetLen)
	copy(result, key)
	paddingText := bytes.Repeat([]byte{byte(left)}, left)
	copy(result[len(key):], paddingText)
	return result
}
