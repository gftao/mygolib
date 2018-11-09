package security

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"mygolib/defs"
	"mygolib/gerror"
)

var IV = []byte{0, 0, 0, 0, 0, 0, 0, 0}
//var IV = []byte("01234567")
func DesEncrypt(origData, key []byte) ([]byte, gerror.IError) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, gerror.New(10000, defs.TRN_SYS_ERROR, err, "des加密失败")
	}
	//origData = ZeroPadding(origData, block.BlockSize())
	origData = PKCS5Padding(origData, block.BlockSize())

	blockMode := cipher.NewCBCEncrypter(block, IV)
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// 3DES加密
func TripleDesEncrypt(origData, key []byte) ([]byte, gerror.IError) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, gerror.New(10010, defs.TRN_SYS_ERROR, err, "3des加密失败")
	}
	origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, IV)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//
func DesDecrypt(crypted, key []byte) ([]byte, gerror.IError) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, gerror.New(10020, defs.TRN_SYS_ERROR, err, "des解密失败")
	}
	blockMode := cipher.NewCBCDecrypter(block, IV)
	//origData := make([]byte, len(crypted))
	origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	//origData = PKCS5UnPadding(origData)

	//origData = ZeroUnPadding(origData)
	return origData, nil
}

// 3DES解密
func TripleDesDecrypt(crypted, key []byte) ([]byte, gerror.IError) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, gerror.New(10030, defs.TRN_SYS_ERROR, err, "3des解密失败")
	}
	blockMode := cipher.NewCBCDecrypter(block, IV)
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = ZeroUnPadding(origData)
	return origData, nil
}

// 2des加密
func DoubleDesEncrypt(origData, key []byte) ([]byte, gerror.IError) {
	a1, err := DesEncrypt(origData, key[:8])
	if err != nil {
		return nil, gerror.New(10040, defs.TRN_SYS_ERROR, err, "2des加密失败")
	}
	a2, err := DesDecrypt(a1, key[8:])
	if err != nil {
		return nil, err
	}
	return DesEncrypt(a2, key[:8])
}

func DoubleCBCDesEncrypt(origData, key []byte) (ciphertext []byte, err gerror.IError) {
	for i := 0; i <= len(origData)-8; i += 8 {
		res, err := DoubleDesEncrypt(origData[i:i+8], key)
		if err != nil {
			return nil, gerror.New(10050, defs.TRN_SYS_ERROR, err, "2des解密失败")
		}
		ciphertext = append(ciphertext, res...)
	}
	return ciphertext, nil
}

func DoubleDesEncryptBase64(origData, key []byte) (string, gerror.IError) {
	res, err := DoubleDesEncrypt(origData, key)
	if err != nil {
		return "", gerror.New(10060, defs.TRN_SYS_ERROR, err, "DoubleDesEncryptBase64 error[%s]", origData)
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

// 2des解密
func DoubleDesDecrypt(crypted, key []byte) ([]byte, gerror.IError) {
	a1, err := DesDecrypt(crypted, key[:8])
	if err != nil {
		return nil, err
	}
	a2, err := DesEncrypt(a1, key[8:])
	if err != nil {
		return nil, err
	}
	return DesDecrypt(a2, key[:8])
}

func DoubleCBCDesDecrypt(ciphertext, key []byte) (origData []byte, err gerror.IError) {
	for i := 0; i <= len(ciphertext)-8; i += 8 {
		res, err := DoubleDesDecrypt(ciphertext[i:i+8], key)
		if err != nil {
			return nil, err
		}
		origData = append(origData, res...)
	}
	return origData, nil
}

func DoubleCBCDesDecryptBase64(crypted string, key []byte) ([]byte, gerror.IError) {
	cipher, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		return nil, gerror.New(10070, defs.TRN_SYS_ERROR, err, "base64.StdEncoding.DecodeString error[%s]", crypted)
	}
	return DoubleCBCDesDecrypt(cipher, key)
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	if len(ciphertext)%blockSize == 0 {
		return ciphertext
	}
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}
