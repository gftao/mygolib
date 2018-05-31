package security

import (
	"crypto/rsa"
	"sort"
	"fmt"
	_"encoding/hex"
	"mygolib/gerror"
	"encoding/base64"
	"crypto"
	"crypto/sha256"
	"mygolib/defs"
)

func UnionRsaSignBase64Url(privSign *rsa.PrivateKey,  keysv map[string]string) (string, gerror.IError) {

	//keysv := make(map[string]string)

	sortKeys := make([]string, 0)
	for key := range keysv {
		sortKeys = append(sortKeys, key)
	}
	sort.Strings(sortKeys)

	var tmpBuf string
	var signBuf string
	for key := range sortKeys {
		tmpBuf = fmt.Sprintf("%s=%s&", sortKeys[key], keysv[sortKeys[key]])
		signBuf += tmpBuf
	}
	signBuf = signBuf[:len(signBuf)-1]

	h := sha256.New()
	h.Write([]byte(signBuf))
	digest := h.Sum(nil)

	/*
	hexBuf := make([]byte, len(digest)*2)
	hex.Encode(hexBuf, digest)
	*/

	ciperdata, err := RsaSignSha256(privSign, digest)
	if err != nil {
		return "", gerror.NewR(1011, err, "RsaSignSha1 error")
	}
	return base64.StdEncoding.EncodeToString(ciperdata), nil
}


func UnionRsaVerify(pubVerify *rsa.PublicKey, respMap map[string]string) (bool,gerror.IError) {

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

	/*h := sha1.New()
	h.Write([]byte(signBuf))
	digest := h.Sum(nil)

	hexBuf := make([]byte, len(digest)*2)
	hex.Encode(hexBuf, digest)*/

/*	signStr, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, gerror.NewR(1022, err, "base64.StdEncoding.DecodeString error[%s]", signature)
	}*/

	ok, err:=RsaVerifySha256Base64(pubVerify, signBuf, signature)
	if err != nil {
		return false,gerror.NewR(1011, err, "RsaSignSha1 error")
			//gerror.NewR(1021, err, "RsaVerifySha1 error")
	}

	if !ok {
		return false,gerror.NewR(1011, err, " 报文验证失败")
	}

	/*err = cert.CheckSignature(x509.SHA1WithRSA, hexBuf, signStr)
	if err != nil {
		return false, gerror.NewR(1021, err, "RsaVerifySha1 error")
	}*/
	return true, nil
}

func RsaSignSha256(privSign *rsa.PrivateKey, data []byte) (ciphertext []byte, err error) {
	/*h := sha256.New()
	h.Write([]byte(data))
	digest := h.Sum(nil)
*/
	ciphertext, err = rsa.SignPKCS1v15(nil, privSign, crypto.SHA256, data)
	if err != nil {
		return nil, gerror.New(10030, defs.TRN_SYS_ERROR, err, "rsa.SignPKCS1v15 error;")
	}
	return
}

func RsaVerifySha256Base64(pubVerify *rsa.PublicKey, data string, signResu string) (ok bool, err error) {

	signVer, err := base64.StdEncoding.DecodeString(signResu)
	if err != nil {
		return false, gerror.New(10060, defs.TRN_SYS_ERROR, err, "base64.StdEncoding.DecodeString error[%s]", signResu)
	}

	h := sha256.New()
	h.Write([]byte(data))
	digest := h.Sum(nil)

	err = rsa.VerifyPKCS1v15(pubVerify, crypto.SHA256, digest, signVer)
	if err != nil {
		return false, gerror.New(10070, defs.TRN_SYS_ERROR, err, "rsa.RsaVerifySha1 error;")
	}
	return true, nil
}