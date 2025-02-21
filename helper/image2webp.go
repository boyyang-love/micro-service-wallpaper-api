package helper

import (
	"bytes"
	"errors"
	"github.com/nickalie/go-webpbin"
)

type CompressedImage struct {
	Buf     *bytes.Buffer
	OriBuf  *bytes.Buffer
	Size    int64
	OriSize int64
}

func Image2Webp(fileBytes *[]byte, quality uint) (compressedImage *CompressedImage, err error) {

	buf := new(bytes.Buffer)
	oriBuf := bytes.NewBuffer(*fileBytes)

	if err = webpbin.NewCWebP().
		Quality(quality).
		Input(bytes.NewReader(*fileBytes)).
		Output(buf).
		Run(); err != nil {
		return nil, errors.New("图片转为webp格式失败")
	}

	reader := bytes.NewReader(buf.Bytes())
	oriRender := bytes.NewReader(oriBuf.Bytes())

	return &CompressedImage{
		Buf:     buf,
		OriBuf:  oriBuf,
		Size:    reader.Size(),
		OriSize: oriRender.Size(),
	}, nil

}
