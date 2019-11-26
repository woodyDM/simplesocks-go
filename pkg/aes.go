package pkg

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func EncryptAsBCB(plainBytes []byte, key []byte, iv []byte) (encrypted []byte) {
	block := getCipher(key)
	blockSize := block.BlockSize()                   // for aes blocksize is always aes.Blocksize
	plainBytes = pkcs5Padding(plainBytes, blockSize) // 补全码
	blockMode := cipher.NewCBCEncrypter(block, iv)   // 加密模式
	encrypted = make([]byte, len(plainBytes))        // 创建数组
	blockMode.CryptBlocks(encrypted, plainBytes)     // 加密
	return encrypted
}

func DecryptAsCBC(encrypted []byte, key []byte, iv []byte) (decrypted []byte) {
	block := getCipher(key)
	blockMode := cipher.NewCBCDecrypter(block, iv)
	decrypted = make([]byte, len(encrypted))    // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted) // 解密
	decrypted = pkcs5UnPadding(decrypted)       // 去除补全码
	return decrypted
}

func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	paddingLen := blockSize - len(cipherText)%blockSize
	paddingText := bytes.Repeat([]byte{byte(paddingLen)}, paddingLen)
	return append(cipherText, paddingText...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPaddingLength := int(origData[length-1])
	return origData[:(length - unPaddingLength)]
}

func getCipher(key []byte) cipher.Block {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return block
}
