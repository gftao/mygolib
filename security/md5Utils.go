package security

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
)

func GenMd5(origData []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write(origData)
	return base64.StdEncoding.EncodeToString(md5Ctx.Sum(nil))
}

func VerifyMd5(origData []byte, desKey string) bool {
	md5Ctx := md5.New()
	md5Ctx.Write(origData)
	return base64.StdEncoding.EncodeToString(md5Ctx.Sum(nil)) == desKey
}

func GetMd5(origData []byte) []byte {
	md5Ctx := md5.New()
	md5Ctx.Write(origData)
	return md5Ctx.Sum(nil)
}

func GetFileMd5(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", (GetMd5(b))), nil
}
