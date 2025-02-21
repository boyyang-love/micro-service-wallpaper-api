package helper

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
)

func MakeFileHash(file multipart.File, fileHeader *multipart.FileHeader) (hash string, err error) {
	h := make([]byte, fileHeader.Size)
	if _, err = file.Read(h); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%x", md5.Sum(h)), nil
	}
}

func MakeImageFileHash(img image.Image, imgType string) (hash string, err error) {
	buf := new(bytes.Buffer)

	switch imgType {
	case "png":
		if err = png.Encode(buf, img); err != nil {
			return "", err
		}
	case "jpeg", "jpg":
		if err = jpeg.Encode(buf, img, nil); err != nil {
			return "", err
		}
	}

	h := make([]byte, buf.Len())
	if _, err = buf.Read(h); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum(h)), nil
}

func MakeImageFileHashByBytes(img []byte) (hash string) {
	return fmt.Sprintf("%x", md5.Sum(img))
}
