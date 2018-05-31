package security_test

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"golib/security"
	"io/ioutil"
	"testing"
	"crypto/sha256"
	"fmt"
	"encoding/base64"
)

func TestRsa(t *testing.T) {
	keyval := "12345678901234567890123456789012"
	workPriKey, err := MyGetRsaPrivateKey("../key_file/a_pri.pem")
	if err != nil {
		t.Error("ParsePKCS1PrivateKey fail:\n", err)
	}
	workPubKey, err := security.GetRsaPublicKey("../key_file/a_pub.pem")
	if err != nil {
		t.Error("GetRsaPublicKey fail:", err)
	}
	res, err := security.RsaEncryptBase64(workPubKey, []byte(keyval))
	if err != nil {
		t.Error("96", "加密失败", err)
		return
	}
	t.Logf(res)
	res = "cKUlch8cHmeUlCcrkO01+aYGWQpZ+De/KqKfmySDwUFJKt32FntoP2zPwIcubasb0oO0TNmX6wszQ76YPoIBxReQ2j2xNTNma3dudSKz1ehWDSektBiml9KaAOb99n5WkfQDtKevm6ZVFasVXMDJKUR1NRSp3nCJMmhQMgIZGYQEYgFVC20XH8YkNft6M/VRyObKGyLi2Fi64J5xZGoXSNB8gadwuVJ7Nodo01/Cc7Qz/zHWoR/BLbm+ejbce8c0ORwcoGBvnbibFnSRBnA+BXvGfaIUa7lUmoLLLwI7C4x34Ss9kKikwQOAhdvLgJYqmIPivSGdBelu7Ay8Dnl/M4hQQp3voLE5xDRGTugcHlcUJsHIbbonIWHCFHR7+sr5c8S+Gbnwi48IrloDzCF09c2XW52+cHqDJ4Dv1p0orIslKiCejqZT4Gpwv74i5NqTUrOMtLwsOCT0uCg/lW01HQ11gF5e1GZtktoVu+MW3P0OnpKVUXfZNnYwk1gCELUwJtcpYsgoDJwg7GEFJiDW5654hsHTSV2LC2vbpQ2xpJwLw6CxObK22w2pAqCBsIrpe9qgVSSjG+/E6ByjlMBrF4p3K3cVKjar7IVOf3aDYCEkhIiiTS7C2/Gc+RYZp7DegSWdvtyscfO9gjeOxszbZ6C/LMcDoUs5tkNCNOok/tg="
	aa, err := security.RsaDecryptBase64(workPriKey, res)
	if err != nil {
		t.Error("96", "解密失败", err)
		return
	}
	t.Log(string(aa))
}

func TestGenRsaPriKey(t *testing.T) {
	pri, pub, err := security.GenRsaKey(2048)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(pri, pub)
}

func MyGetRsaPrivateKey(filepath string) (pubKey *rsa.PrivateKey, err error) {
	keyBuf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, errors.New("1")
	}
	block, _ := pem.Decode(keyBuf)
	if block == nil {
		return nil, errors.New("2")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("3")
	}
	return priv.(*rsa.PrivateKey), nil

}

func TestRsa2(t *testing.T)  {
	//pubKey := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCEewQr1rdY53muMTPS1ekmmdIfvk08Udv+swQp80D3kZVgt9+CLBn+2AJ7d/P7T8xK/do2lEpjcATzMDE29cUkRn1zZMWnMwrRTEtfMgwIWFp3JhgYBuu3jFqAcl7eIlU5yHAwnlgnGApWqscBh/Qt4A44xO0TXX0zhYuAPqxUGwIDAQAB"
//	priKey := `-----BEGIN RSA PRIVATE KEY-----
//MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAIR7BCvWt1jnea4xM9LV6SaZ0h++TTxR2/6zBCnzQPeRlWC334IsGf7YAnt38/tPzEr92jaUSmNwBPMwMTb1xSRGfXNkxaczCtFMS18yDAhYWncmGBgG67eMWoByXt4iVTnIcDCeWCcYClaqxwGH9C3gDjjE7RNdfTOFi4A+rFQbAgMBAAECgYA27NWe40dimZ1uJcIJoFwof8+JD5nv7zRZVZjdV5fQzE/1KGaHDoe8i4wD6oiB4eSeFr74r+Rjc5bpyEovMhgIRFH+St67rz/Rxvr3R75M11ef/rtyrTdzHYJAUDfk3E/j4UJVOMvk+2izDwh9eDIioF5oC8Fr504tXiKY1ARm8QJBAOf/wXDizyiKbyvxUyDbE9XQyCR7dlreaWXZh498aarK0VofYufkYxZoxofxo8FY0xbLHozEOTzQ2V3Eivvjj6kCQQCSL57H1FkSp/H6/ycI/SExBn28fUDIBi3j+RLeAIwmxRNjZM0xeynyjyhniqMfj9fFau8VtYNzFcagMDDKoDAjAkEA2CA1mDFjJXRZfslRVNFimBTo7ruplZuO+pf8ppoTYk2RXHDS1g64lH7FPI3KrOtPsvNEoYSHgfVaGfVoOKJVCQJATijn3C/M8AybdHe3hzbP6EZwM7dES64CG0GwtMHWLWRxWVMr4qjXZLjmAXY+gUGHPCZQbmr+PSoHaN5bN/stwQJAJIN76UUKTsLXGtfYkPQp7OLT7X9aphs6dHLFYYChWQBUEQ4Znvt3WplhATnbKmx+1DymAqH8FQWNMJAD3/+nBA==
//-----END RSA PRIVATE KEY-----`
	pri, err := security.GetRsaPrivateKey("./pri_pa.key")
	if err != nil {
		t.Fatal(err)
	}
	//pub, err := security.GetRsaPublicKeyByString(pubKey)
	//if err != nil {
	//	t.Fatal(err)
	//}
	data := "dfasdfsadgadfsfasdfasdfsdfadsfgdsfasdfasfdas3214"
	h := sha256.New()
	h.Write([]byte(data))
	digest := h.Sum(nil)
	ciperdata, err := security.RsaSignSha256(pri, digest)
	fmt.Println(base64.StdEncoding.EncodeToString(ciperdata))

}