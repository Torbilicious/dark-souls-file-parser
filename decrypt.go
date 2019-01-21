package main

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func decrypt(key []byte, secure []byte) (decoded []byte) {
	cipherText := make([]byte, len(secure))

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	//decodedmess = string(cipherText)
	return cipherText
}

//func decrypt(key []byte, secure []byte) (decoded []byte) {
//	bc, _ := aes.NewCipher(key)
//
//	decoded = make([]byte, 0)
//	for i := 0; i< len(secure);i+=aes.BlockSize  {
//		dst := make([]byte, len(secure))
//		bc.Decrypt(dst, secure)
//
//		for _, oneByte := range dst {
//			decoded = append(decoded, oneByte)
//		}
//	}
//
//
//	return decoded
//}
