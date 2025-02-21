package helper

import (
	"crypto/md5"
	"fmt"
	"io"
)

func MakeHash(s string) (hs string, err error) {
	h := md5.New()
	if _, err = io.WriteString(h, s); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%x", h.Sum(nil)), err
	}
}
