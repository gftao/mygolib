package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)


type SvrCert struct {
	KeyType            string `json:"key_type"`
	KeyFile            string `json:"key_file"`
	CertFile           string `json:"cert_file"`
	SerialNumber       string `json:"serial_number"`
	X509Cert           *x509.Certificate
	PrivateKey         interface{}
	PublicKey          interface{}
}

func NewSvrCert(jsonConfig string) (*SvrCert, error) {

	cert := new(SvrCert)
	err := json.Unmarshal([]byte(jsonConfig), cert)
	if err != nil {
		return nil, err
	}

	keyBuf, err := ioutil.ReadFile(cert.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("私钥文件[%s]读取失败[%s]", cert.KeyFile, err)
	}

	keyBlock, _ := pem.Decode(keyBuf)
	if keyBlock == nil {
		return nil, fmt.Errorf("私钥文件[%s]PEM解码失败", cert.KeyFile)
	}

	certBuf, err:= ioutil.ReadFile(cert.CertFile)
	if err != nil {
		return nil, fmt.Errorf("证书文件[%s]读取失败[%s]", cert.CertFile, err)
	}

	certBlock, _ := pem.Decode(certBuf)
	if keyBlock == nil {
		return nil, fmt.Errorf("证书文件[%s]PEM解码失败", cert.CertFile)
	}

	switch cert.KeyType {
	case FILE_RSA_PEM_PRIV:
		priv, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("解析PEM私钥失败[%s]", err)
		}
		cert.PrivateKey = priv
		cert.PublicKey = priv.Public()
		cert.SerialNumber = ""
	case FILE_CERT_PFX:
		fallthrough
	default:
		priv, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("解析PEM私钥失败[%s]", err)
		}
		cert.X509Cert, err = x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("解析PEM证书[%s]失败[%s]",
				cert.CertFile, err)
		}
		cert.PrivateKey = priv
		cert.PublicKey = priv.Public()
		cert.SerialNumber = cert.X509Cert.SerialNumber.String()

	}
	return cert, nil
}

/*algo 算法:
  signbuf 签名串
  signature 签名值
*/
func (cert *SvrCert) Sign(algo int, signbuf []byte) ([]byte, error) {

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

	switch priv := cert.PrivateKey.(type) {
	case *rsa.PrivateKey:
		return priv.Sign(rand.Reader, digest, hashType)
	default:
		return nil, x509.ErrUnsupportedAlgorithm
	}
}

/*
	公钥加密:
*/
func (cert *SvrCert) EncryptPKCS1v15(origdata []byte) ([]byte, error) {
	switch pub := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		rsaKey, err := LoadFromRSAPublic(pub)
		if err != nil {
			return nil, err
		}
		return rsaKey.EncryptPKCS1v15(origdata)
	default:
		return nil, x509.ErrUnsupportedAlgorithm
	}
}

/*

	私钥解密
*/
func (cert *SvrCert) DecryptPKCS1v15(cipherdata []byte) ([]byte, error) {
	switch priv := cert.PrivateKey.(type) {
	case *rsa.PrivateKey:
		rsaKey, err := LoadFromRSAPrivate(priv)
		if err != nil {
			return nil, err
		}
		return rsaKey.DecryptPKCS1v15(cipherdata)
	default:
		return nil, x509.ErrUnsupportedAlgorithm
	}
}

/*获取证书信息*/
func (cert *SvrCert) GetSerialNumber() string {
	return cert.SerialNumber
}
