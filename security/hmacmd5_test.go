package security

import (
	"testing"
)

const (
	hmaced string = "812BA54471064FAC6C19E98EEDA2D1BB"
	src    string = "1234567890"
	key    string = "1234567890"
)

func TestHmacMd5(t *testing.T) {
	t.Parallel()
	if HmacMd5(src, key) != hmaced {
		t.FailNow()
	}
}

func TestVerifyHmacMd5(t *testing.T) {
	t.Parallel()
	if VerifyHmacMd5([]byte(src), []byte(hmaced), []byte(key)) {
		t.FailNow()
	}
}
