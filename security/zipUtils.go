package security

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"io"
)

func Zlib(src []byte) []byte {
	buf := bytes.Buffer{}
	w := zlib.NewWriter(&buf)
	w.Write(src)
	w.Close()
	return buf.Bytes()
}

func UnZlib(src []byte) ([]byte, error) {
	buf := bytes.Buffer{}
	br := bytes.NewReader(src)
	r, err := zlib.NewReader(br)
	if err != nil {
		return nil, err
	}
	io.Copy(&buf, r)
	return buf.Bytes(), nil
}

func ZlibBase64(src []byte) string {
	return base64.StdEncoding.EncodeToString(Zlib(src))
}

func UnZlibBase64(src string) ([]byte, error) {
	des, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}
	return UnZlib(des)
}
