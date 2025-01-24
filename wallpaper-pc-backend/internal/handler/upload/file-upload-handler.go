package upload

import (
	"bytes"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/logic/upload"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/uploadclient"
	"github.com/zeromicro/go-zero/rest/httpx"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
)

func FileUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileUploadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		// 用户上传名称 以及路径（blog,image）
		fileCustomName := r.PostFormValue("file_name")
		fileCustomFolder := r.PostFormValue("folder")

		file, err = fileHeader.Open()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		defer func(file multipart.File) {
			_ = file.Close()
		}(file)

		img, imgType, err := image.Decode(file)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		originBuffer := new(bytes.Buffer)
		switch imgType {
		case "png":
			if err = png.Encode(originBuffer, img); err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}
		case "jpeg", "jpg":
			if err = jpeg.Encode(originBuffer, img, nil); err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}
		}

		l := upload.NewFileUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileUpload(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			fmt.Println(originBuffer.Bytes())
			fileUpload, err := svcCtx.UploadService.FileUpload(
				r.Context(),
				&uploadclient.FileUploadReq{
					File:       originBuffer.Bytes(),
					FileName:   fileCustomName,
					OriFolder:  fmt.Sprintf("BOYYANG/%d/%s/%s", 555, fileCustomFolder, "origin"),
					Folder:     fmt.Sprintf("BOYYANG/%d/%s/%s", 555, fileCustomFolder, "compressed"),
					FileType:   imgType,
					W:          uint32(img.Bounds().Dx()),
					H:          uint32(img.Bounds().Dy()),
					Quality:    100,
					UserId:     "555",
					BucketName: "boyyang",
					Type:       req.Type,
				})

			if err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}

			resp.Data = types.FileUploadResdata{
				FileName:   fileCustomName,
				Path:       fileUpload.FilePath,
				OriginPath: fileUpload.OriFilePath,
			}

			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
