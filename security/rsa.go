package security

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

type RSAKey struct {
	KeyType    string `json:"key_type"`
	KeyFile    string `json:"key_file"`
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
	Modules    int
}

func NewRSA(jsonConfig string) (*RSAKey, error) {
	r := new(RSAKey)
	err := json.Unmarshal([]byte(jsonConfig), r)
	if err != nil {
		return nil, err
	}

	keyBuf, err := ioutil.ReadFile(r.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("读取RSA密钥失败[%s]", r.KeyFile)
	}

	block, _ := pem.Decode(keyBuf)
	if len(block.Bytes) == 0 {
		return nil, fmt.Errorf("读取RSA密钥文件w为空[%s]", r.KeyFile)
	}

	var ok bool
	switch r.KeyType {
	case FILE_RSA_PEM_PUB:
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("解析RSA公钥失败[%s]", err)
		}
		r.PublicKey, ok = pub.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("Value returned from ParsePKIXPublicKey was not an RSA public key")
		}
		r.Modules = (r.PublicKey.N.BitLen() + 7) / 8
	case FILE_RSA_PEM_PRIV:
		r.PrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("ParsePKCS1PrivateKey err:%s", err)
		}
		r.Modules = (r.PrivateKey.N.BitLen() + 7) / 8
	default:
		return nil, fmt.Errorf("非法的RSA公钥类型[%s]", r.KeyType)
	}
	return r, nil
}

/*使用RSA公钥装载*/
func LoadFromRSAPublic(pub *rsa.PublicKey) (*RSAKey, error) {
	r := new(RSAKey)
	r.PublicKey = pub
	r.Modules = (r.PublicKey.N.BitLen() + 7) / 8
	return r, nil
}

/*使用RSA私钥装载*/
func LoadFromRSAPrivate(priv *rsa.PrivateKey) (*RSAKey, error) {
	r := new(RSAKey)
	r.PrivateKey = priv
	r.Modules = (r.PrivateKey.N.BitLen() + 7) / 8
	return r, nil
}

func (r *RSAKey) Sign(algo int, signbuf []byte) ([]byte, error) {
	var hashType crypto.Hash

	switch algo {
	case SHA1WithRSA, DSAWithSHA1, ECDSAWithSHA1:
		hashType = crypto.SHA1
	case SHA256WithRSA, DSAWithSHA256, ECDSAWithSHA256:
		hashType = crypto.SHA256
	case SHA384WithRSA, ECDSAWithSHA384:
		hashType = crypto.SHA384
	case SHA512WithRSA, ECDSAWithSHA512:
		hashType = crypto.SHA512
	case MD2WithRSA, MD5WithRSA:
		return nil, x509.InsecureAlgorithmError(algo)
	default:
		return nil, fmt.Errorf("非法的签名算法:[%d]", algo)
	}

	if !hashType.Available() {
		return nil, x509.ErrUnsupportedAlgorithm
	}
	h := hashType.New()
	h.Write(signbuf)
	digest := h.Sum(nil)

	return r.PrivateKey.Sign(rand.Reader, digest, hashType)
}

func (r *RSAKey) Verfy(algo int, signbuf []byte, signature []byte) error {
	var hashType crypto.Hash

	switch algo {
	case SHA1WithRSA, DSAWithSHA1, ECDSAWithSHA1:
		hashType = crypto.SHA1
	case SHA256WithRSA, DSAWithSHA256, ECDSAWithSHA256:
		hashType = crypto.SHA256
	case SHA384WithRSA, ECDSAWithSHA384:
		hashType = crypto.SHA384
	case SHA512WithRSA, ECDSAWithSHA512:
		hashType = crypto.SHA512
	case MD2WithRSA, MD5WithRSA:
		return x509.InsecureAlgorithmError(algo)
	default:
		return fmt.Errorf("非法的验签算法:[%d]", algo)
	}

	if !hashType.Available() {
		return x509.ErrUnsupportedAlgorithm
	}
	h := hashType.New()
	h.Write(signbuf)
	digest := h.Sum(nil)

	return rsa.VerifyPKCS1v15(r.PublicKey, hashType, digest, signature)

}

func (r *RSAKey) EncryptPKCS1v15(origdata []byte) ([]byte, error) {
	encLen := r.Modules - 11 //PKCS1v15
	orig := bytes.NewBuffer(origdata)

	var ciperdata bytes.Buffer

	for {
		data := orig.Next(encLen)
		if len(data) == 0 {
			break
		}
		ciper, err := rsa.EncryptPKCS1v15(rand.Reader, r.PublicKey, data)
		if err != nil {
			return nil, fmt.Errorf("rsa.EncryptPKCS1v15 error[%s]", err)
		}
		ciperdata.Write(ciper)
	}
	return ciperdata.Bytes(), nil
}

func (r *RSAKey) DecryptPKCS1v15(cipherdata []byte) ([]byte, error) {
	encLen := r.Modules //PKCS1v15
	orig := bytes.NewBuffer(cipherdata)

	var origdata bytes.Buffer

	for {
		data := orig.Next(encLen)
		if len(data) == 0 {
			break
		}
		ciper, err := rsa.DecryptPKCS1v15(rand.Reader, r.PrivateKey, data)
		if err != nil {
			return nil, fmt.Errorf("rsa.DecryptPKCS1v15 error[%s]", err)
		}
		origdata.Write(ciper)
	}
	return origdata.Bytes(), nil

}

func (r *RSAKey) GetModules() int {
	return r.Modules
}
