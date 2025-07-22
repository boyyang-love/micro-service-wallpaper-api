package upload

import (
	"bytes"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-api/helper"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	upload2 "github.com/boyyang-love/micro-service-wallpaper-rpc/upload/pb/upload"
	"github.com/zeromicro/go-zero/rest/httpx"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"net/http"
	"strings"
)

func ImageUploadByBytesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ImageUploadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		userid := fmt.Sprintf("%s", r.Context().Value("Id"))

		file := r.FormValue("file")

		img, imgType, err := image.Decode(bytes.NewBufferString(file))
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

		hash := helper.MakeImageFileHashByBytes(originBuffer.Bytes())

		oriPath := fmt.Sprintf("%s/%s/%s/%s/%s%s", req.RootDir, req.Dir, userid, "original", hash, fmt.Sprintf(".%s", imgType))
		comPath := fmt.Sprintf("%s/%s/%s/%s/%s%s", req.RootDir, req.Dir, userid, "compress", hash, ".webp")

		is, info := Is(svcCtx.DB, hash, oriPath)

		if is {
			httpx.OkJsonCtx(r.Context(), w, types.ImageUploadRes{
				Base: types.Base{
					Code: 1,
					Msg:  "文件上传成功",
				},
				Data: types.ImageUploadResdata{
					Id:         info.Id,
					FileName:   info.FileName,
					Path:       info.FilePath,
					OriginPath: info.OriginFilePath,
				},
			})
			return
		} else {
			imageUpload, err := svcCtx.
				UploadService.
				ImageUpload(
					r.Context(),
					&upload2.ImageUploadReq{
						File:       originBuffer.Bytes(),
						Path:       comPath,
						OriPath:    oriPath,
						Quality:    req.Quality,
						BucketName: req.BucketName,
					})
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}

			uploadInfo := models.Upload{
				Hash:           imageUpload.Data.OriETag,
				FileName:       req.FileName,
				OriginFileSize: int64(imageUpload.Data.OriSize),
				FileSize:       int64(imageUpload.Data.Size),
				OriginType:     imgType,
				FileType:       "webp",
				OriginFilePath: oriPath,
				FilePath:       comPath,
				Type:           req.Type,
				W:              img.Bounds().Dx(),
				H:              img.Bounds().Dy(),
				Status:         req.Status,
				UserId:         userid,
			}
			if err := svcCtx.
				DB.
				Model(&models.Upload{}).
				Create(&uploadInfo).Error; err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}

			if req.Tags != "" {
				var uploadTags []models.UploadTag

				for _, v := range strings.Split(req.Tags, ",") {
					uploadTags = append(
						uploadTags,
						models.UploadTag{
							UploadId: uploadInfo.Id,
							TagId:    v,
						},
					)
				}

				if err = svcCtx.DB.
					Model(&models.UploadTag{}).
					Create(&uploadTags).
					Error; err != nil {
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
			}

			if req.Category != "" {
				var uploadCategory []models.UploadCategory

				for _, v := range strings.Split(req.Category, ",") {
					uploadCategory = append(
						uploadCategory,
						models.UploadCategory{
							UploadId:   uploadInfo.Id,
							CategoryId: v,
						},
					)
				}

				if err = svcCtx.DB.
					Model(&models.UploadCategory{}).
					Create(&uploadCategory).
					Error; err != nil {
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
			}

			if req.Recommend != "" {
				var uploadRecommend []models.UploadRecommend

				for _, v := range strings.Split(req.Recommend, ",") {
					uploadRecommend = append(
						uploadRecommend,
						models.UploadRecommend{
							UploadId:    uploadInfo.Id,
							RecommendId: v,
						},
					)
				}

				if err = svcCtx.DB.
					Model(&models.UploadRecommend{}).
					Create(&uploadRecommend).
					Error; err != nil {
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
			}

			if req.Group != "" {
				var uploadGroup []models.UploadGroup
				for _, v := range strings.Split(req.Group, ",") {
					uploadGroup = append(uploadGroup, models.UploadGroup{
						UploadId: uploadInfo.Id,
						GroupId:  v,
					})
				}

				if err = svcCtx.DB.
					Model(&models.UploadGroup{}).
					Create(&uploadGroup).
					Error; err != nil {
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
			}

			httpx.OkJsonCtx(r.Context(), w, &types.ImageUploadRes{
				Base: types.Base{
					Code: 1,
					Msg:  "文件上传成功",
				},
				Data: types.ImageUploadResdata{
					Id:         uploadInfo.Id,
					FileName:   req.FileName,
					Path:       imageUpload.Data.Path,
					OriginPath: imageUpload.Data.OriPath,
				},
			})
			return
		}
	}
}
