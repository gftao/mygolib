package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"mygolib/defs"
	e "mygolib/gerror"
	"io/ioutil"
	"sort"
)

func GenRsaKey(bits int) (priRes, pubRes string, err error) {
	priKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	derStream := x509.MarshalPKCS1PrivateKey(priKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	priRes = string(pem.EncodeToMemory(block))

	//生成公钥
	pubKey := &priKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", "", nil
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubRes = string(pem.EncodeToMemory(block))

	return
}

func RsaSignBase64(privSign *rsa.PrivateKey, req map[string]string) (ciphertext string, err error) {

	sortKeys := make([]string, 0)
	for key := range req {
		sortKeys = append(sortKeys, key)
	}
	sort.Strings(sortKeys)

	var tmpBuf string
	var signBuf string
	for key := range sortKeys {
		tmpBuf = fmt.Sprintf("%s=%s&", sortKeys[key], req[sortKeys[key]])
		signBuf += tmpBuf
	}
	signBuf = signBuf[:len(signBuf)-1]

	h := sha1.New()
	h.Write([]byte(signBuf))
	digest := h.Sum(nil)

	hexBuf := make([]byte, len(digest)*2)
	hex.Encode(hexBuf, digest)

	ciperdata, err := RsaSignSha1(privSign, hexBuf)
	if err != nil {
		return "", e.New(10000, defs.TRN_SYS_ERROR, err, "RsaSignSha1 error")
	}
	return base64.StdEncoding.EncodeToString(ciperdata), nil
}

func RsaVerifyBase64(cert *x509.Certificate, respMap map[string]string) (ok bool, err error) {

	var sortKeys []string
	var signature string
	var tmpBuf string
	var signBuf string

	for key := range respMap {
		if key == "signature" {
			signature = respMap[key] // 签名结果
		} else {
			sortKeys = append(sortKeys, key)
		}

	}
	sort.Strings(sortKeys)

	for key := range sortKeys {
		tmpBuf = fmt.Sprintf("%s=%s&", sortKeys[key], respMap[sortKeys[key]])
		signBuf += tmpBuf
	}
	signBuf = signBuf[:len(signBuf)-1]

	h := sha1.New()
	h.Write([]byte(signBuf))
	digest := h.Sum(nil)

	hexBuf := make([]byte, len(digest)*2)
	hex.Encode(hexBuf, digest)

	signStr, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, e.New(10010, defs.TRN_SYS_ERROR, err, "base64.StdEncoding.DecodeString error[%s]", signature)
	}

	err = cert.CheckSignature(x509.SHA1WithRSA, hexBuf, signStr)
	if err != nil {
		return false, e.New(10020, defs.TRN_SYS_ERROR, err, "RsaVerifySha1 error")
	}
	return true, nil
}

func RsaSignSha1(privSign *rsa.PrivateKey, data []byte) (ciphertext []byte, err error) {
	h := sha1.New()
	h.Write([]byte(data))
	digest := h.Sum(nil)

	ciphertext, err = rsa.SignPKCS1v15(nil, privSign, crypto.SHA1, digest)
	if err != nil {
		return nil, e.New(10030, defs.TRN_SYS_ERROR, err, "rsa.SignPKCS1v15 error;")
	}
	return
}

func RsaSignSha1Base64(privSign *rsa.PrivateKey, data []byte) (string, error) {
	h := sha1.New()
	h.Write([]byte(data))
	digest := h.Sum(nil)

	ciphertext, err := rsa.SignPKCS1v15(nil, privSign, crypto.SHA1, digest)
	if err != nil {
		return "", e.New(10040, defs.TRN_SYS_ERROR, err, "rsa.SignPKCS1v15 error;")
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func RsaVerifySha1(pubVerify *rsa.PublicKey, data []byte, signResu []byte) (ok bool, err error) {
	h := sha1.New()
	h.Write([]byte(data))
	digest := h.Sum(nil)

	err = rsa.VerifyPKCS1v15(pubVerify, crypto.SHA1, digest, signResu)
	if err != nil {
		return false, e.New(10050, defs.TRN_SYS_ERROR, err, "rsa.RsaVerifySha1 error;")
	}
	return true, nil
}

func RsaVerifySha1Base64(pubVerify *rsa.PublicKey, data string, signResu string) (ok bool, err error) {

	signVer, err := base64.StdEncoding.DecodeString(signResu)
	if err != nil {
		return false, e.New(10060, defs.TRN_SYS_ERROR, err, "base64.StdEncoding.DecodeString error[%s]", signResu)
	}

	h := sha1.New()
	h.Write([]byte(data))
	digest := h.Sum(nil)

	err = rsa.VerifyPKCS1v15(pubVerify, crypto.SHA1, digest, signVer)
	if err != nil {
		return false, e.New(10070, defs.TRN_SYS_ERROR, err, "rsa.RsaVerifySha1 error;")
	}
	return true, nil
}

//获取公钥
func GetRsaPublicKey(filePath string) (pubKey *rsa.PublicKey, err error) {
	keyBuf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, e.New(10080, defs.TRN_SYS_ERROR, err, "ReadFile[%s] err", filePath)
	}

	block, _ := pem.Decode(keyBuf)
	if block == nil {
		return nil, e.New(10090, defs.TRN_SYS_ERROR, err, "pem.Decode[%s] err", keyBuf)
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, e.New(10100, defs.TRN_SYS_ERROR, err, "ReadFile[%s] err", filePath)
	}

	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, e.New(10110, defs.TRN_SYS_ERROR, nil, "pubInterface.(*rsa.PublicKey) error")
	}

	return pubKey, nil

}
func GetRsaPublicKeyByString(pubkey string) (pubKey *rsa.PublicKey, err error) {

	block, _ := pem.Decode([]byte(pubkey))
	if block == nil {
		return nil, e.New(10090, defs.TRN_SYS_ERROR, err, "pem.Decode[%s] err", pubkey)
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, e.New(10100, defs.TRN_SYS_ERROR, err, "ReadString[%s] err", pubkey)
	}

	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, e.New(10110, defs.TRN_SYS_ERROR, nil, "pubInterface.(*rsa.PublicKey) error")
	}

	return pubKey, nil

}

// 获取私钥
func GetRsaPrivateKey(filepath string) (pubKey *rsa.PrivateKey, err error) {
	keyBuf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, e.New(10120, defs.TRN_SYS_ERROR, err, "ReadFile[%s] err", filepath)
	}
	block, _ := pem.Decode(keyBuf)
	if block == nil {
		return nil, e.New(10130, defs.TRN_SYS_ERROR, err, "pem.Decode[%s] err", keyBuf)
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, e.New(10140, defs.TRN_SYS_ERROR, err, "ParsePKCS1PrivateKey err")
	}
	return priv, nil

}
func GetRsaPrivatePKCS8Key(filepath string) (pubKey *rsa.PrivateKey, err error) {
	keyBuf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, e.New(10120, defs.TRN_SYS_ERROR, err, "ReadFile[%s] err", filepath)
	}
	block, _ := pem.Decode(keyBuf)
	if block == nil {
		return nil, e.New(10130, defs.TRN_SYS_ERROR, err, "pem.Decode[%s] err", keyBuf)
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, e.New(10140, defs.TRN_SYS_ERROR, err, "ParsePKCS1PrivateKey err")
	}
	return priv.(*rsa.PrivateKey), nil

}
func GetRsaPrivateKeyByString(prikey string) (pubKey *rsa.PrivateKey, err error) {
	block, _ := pem.Decode([]byte(prikey))
	if block == nil {
		return nil, e.New(10130, defs.TRN_SYS_ERROR, err, "pem.Decode[%s] err", prikey)
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, e.New(10140, defs.TRN_SYS_ERROR, err, "ParsePKCS1PrivateKey err")
	}
	return priv, nil

}

/*RSA 报文加密
 */
func RsaEncryptBase64(encKey *rsa.PublicKey, origData []byte) (ciphertext string, err error) {

	maxSize := RsaPubKeyModules(encKey) - 11 //PKCS1v15
	inputLen := len(origData)
	encLen := 0
	var ciperdata []byte
	for offset := 0; offset < inputLen; offset += maxSize {
		if inputLen-offset >= maxSize {
			encLen = maxSize
		} else {
			encLen = inputLen - offset
		}

		ciper, err := rsa.EncryptPKCS1v15(rand.Reader, encKey, origData[offset:offset+encLen])
		if err != nil {
			return "", e.New(10150, defs.TRN_SYS_ERROR, err, "rsa.EncryptPKCS1v15 error")
		}
		ciperdata = append(ciperdata, ciper...)
	}

	ciphertext = base64.StdEncoding.EncodeToString(ciperdata)
	return ciphertext, nil
}

// 解密
func RsaDecryptBase64(decKey *rsa.PrivateKey, ciphertext string) (origData []byte, err error) {

	maxSize := RsaPrivKeyModules(decKey) //PKCS1v15

	cipher, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, e.New(10160, defs.TRN_SYS_ERROR, err, "base64.StdEncoding.DecodeString error[%s]", ciphertext)
	}

	cipherLen := len(cipher)
	decLen := 0
	for offset := 0; offset < cipherLen; offset += maxSize {
		if cipherLen-offset >= maxSize {
			decLen = maxSize
		} else {
			decLen = cipherLen - offset
		}
		data, err := rsa.DecryptPKCS1v15(rand.Reader, decKey, cipher[offset:offset+decLen])
		if err != nil {
			return nil, e.New(10170, defs.TRN_SYS_ERROR, err, "rsa.DecryptPKCS1v15 error")
		}
		origData = append(origData, data...)
	}

	return origData, nil
}

func RsaPubKeyModules(pub *rsa.PublicKey) int {
	return (pub.N.BitLen() + 7) / 8
}

func RsaPrivKeyModules(priv *rsa.PrivateKey) int {
	return (priv.N.BitLen() + 7) / 8
}
