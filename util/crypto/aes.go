package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// AESEncryptCBC encrypt plainText with CBC mode, the key length can be 16, 24, or 32 bytes
func AESEncryptCBC(plainText, key []byte) ([]byte, error) {
	if len(plainText) == 0 || len(key) == 0 {
		return nil, errors.New("param is invalid, plainText and key can not be empty")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	paddedPlainText := PKCS7Padding(plainText, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	cipherByte := make([]byte, len(paddedPlainText))
	// 加密
	blockMode.CryptBlocks(cipherByte, paddedPlainText)
	return cipherByte, nil
}

// AESDecryptCBC decrypt cipher text with CBC mode，the key length can be 16, 24, or 32 bytes
func AESDecryptCBC(cipherByte, key []byte) ([]byte, error) {
	if len(cipherByte) == 0 || len(key) == 0 {
		return nil, errors.New("param is invalid, cipherByte and key can not be empty")
	}
	//k := []byte(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	if len(cipherByte)%blockSize != 0 {
		return nil, errors.New("cipher text length is invalid")
	}
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	plainTextByte := make([]byte, len(cipherByte))
	// 解密
	blockMode.CryptBlocks(plainTextByte, cipherByte)
	// 去补全码
	plainTextByte, err = PKCS7UnPadding(plainTextByte)
	return plainTextByte, err
}

func PKCS7Padding(cipherText []byte, blocksize int) []byte {
	padding := blocksize - len(cipherText)%blocksize
	paddedtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, paddedtext...)
}

func PKCS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length <= 0 {
		return nil, errors.New("cipher text is invalid, unpadded failed, data length is zero")
	}
	unpadded := int(origData[length-1])
	if length-unpadded < 0 {
		return nil, errors.New("cipher text is invalid, unpadded failed")
	}
	return origData[:(length - unpadded)], nil
}
