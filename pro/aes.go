package pro

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"log"
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

// =================== CFB ======================

func EncryptAsCFB(plainByte []byte, key []byte, iv []byte) (encrypted []byte) {
	block := getCipher(key)
	encrypted = make([]byte, aes.BlockSize+len(plainByte))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], plainByte)
	return encrypted
}
func DecryptAsCFB(encrypted []byte, key []byte, iv []byte) (decrypted []byte) {
	log.Printf("aes dec : raw size is %d\n", len(encrypted))
	block := getCipher(key)
	if len(encrypted) < aes.BlockSize {
		panic("cipherText too short, should greater than aes blockSize 16. ")
	}
	encrypted = encrypted[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}

func getCipher(key []byte) cipher.Block {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return block
}
