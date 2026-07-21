package upload

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"
	"sync"

	"github.com/boyyang-love/micro-service-wallpaper-api/helper"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	upload2 "github.com/boyyang-love/micro-service-wallpaper-rpc/upload/pb/upload"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var (
	uploadSemaphore chan struct{}
	once            sync.Once
)

func initSemaphore(maxConcurrent int) {
	once.Do(func() {
		if maxConcurrent <= 0 {
			maxConcurrent = 10
		}
		uploadSemaphore = make(chan struct{}, maxConcurrent)
	})
}

func ImageUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initSemaphore(svcCtx.Config.UploadConf.MaxConcurrent)
		uploadSemaphore <- struct{}{}
		defer func() { <-uploadSemaphore }()

		var req types.ImageUploadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		userid := fmt.Sprintf("%s", r.Context().Value("Id"))

		maxFileSize := svcCtx.Config.UploadConf.MaxFileSize
		if maxFileSize <= 0 {
			maxFileSize = 10 << 20
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			if err.Error() == "http: request body too large" {
				httpx.ErrorCtx(r.Context(), w, fmt.Errorf("文件大小超过限制，最大允许 %dMB", maxFileSize/(1<<20)))
				return
			}
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		defer file.Close()

		logx.Infof("用户 %s 上传文件: %s, 大小: %d bytes", userid, fileHeader.Filename, fileHeader.Size)

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

		hash := helper.MakeImageFileHashByBytes(originBuffer.Bytes())
		ext := helper.FileNameExt(fileHeader.Filename)

		oriPath := fmt.Sprintf("%s/%s/%s/%s/%s%s", req.RootDir, req.Dir, userid, "original", hash, ext)
		comPath := fmt.Sprintf("%s/%s/%s/%s/%s%s", req.RootDir, req.Dir, userid, "compress", hash, ".webp")

		is, info := Is(svcCtx.DB, hash, oriPath)

		if is {
			httpx.OkJsonCtx(r.Context(), w, types.ImageUploadRes{
				Base: types.Base{Code: 1, Msg: "文件上传成功"},
				Data: types.ImageUploadResdata{
					Id:         info.Id,
					FileName:   info.FileName,
					Path:       info.FilePath,
					OriginPath: info.OriginFilePath,
				},
			})
			return
		} else {
			imageUpload, err := svcCtx.UploadService.CosUpload(
				r.Context(),
				&upload2.ImageUploadReq{
					File:       originBuffer.Bytes(),
					Path:       comPath,
					OriPath:    oriPath,
					Quality:    req.Quality,
					BucketName: req.BucketName,
				},
			)
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
			if err := svcCtx.DB.Model(&models.Upload{}).Create(&uploadInfo).Error; err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}

			if req.Tags != "" {
				var uploadTags []models.UploadTag
				for _, v := range strings.Split(req.Tags, ",") {
					uploadTags = append(uploadTags, models.UploadTag{UploadId: uploadInfo.Id, TagId: v})
				}
				if err = svcCtx.DB.Model(&models.UploadTag{}).Create(&uploadTags).Error; err != nil {
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
			}

			if req.Category != "" {
				var uploadCategory []models.UploadCategory
				for _, v := range strings.Split(req.Category, ",") {
					uploadCategory = append(uploadCategory, models.UploadCategory{UploadId: uploadInfo.Id, CategoryId: v})
				}
				if err = svcCtx.DB.Model(&models.UploadCategory{}).Create(&uploadCategory).Error; err != nil {
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
			}

			if req.Recommend != "" {
				var uploadRecommend []models.UploadRecommend
				for _, v := range strings.Split(req.Recommend, ",") {
					uploadRecommend = append(uploadRecommend, models.UploadRecommend{UploadId: uploadInfo.Id, RecommendId: v})
				}
				if err = svcCtx.DB.Model(&models.UploadRecommend{}).Create(&uploadRecommend).Error; err != nil {
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
			}

			if req.Group != "" {
				var uploadGroup []models.UploadGroup
				for _, v := range strings.Split(req.Group, ",") {
					uploadGroup = append(uploadGroup, models.UploadGroup{UploadId: uploadInfo.Id, GroupId: v})
				}
				if err = svcCtx.DB.Model(&models.UploadGroup{}).Create(&uploadGroup).Error; err != nil {
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
			}

			if req.Album != "" {
				var uploadAlbum []models.UploadAlbum
				for _, v := range strings.Split(req.Album, ",") {
					uploadAlbum = append(uploadAlbum, models.UploadAlbum{UploadId: uploadInfo.Id, AlbumId: v})
				}
				if err = svcCtx.DB.Model(&models.UploadAlbum{}).Create(&uploadAlbum).Error; err != nil {
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
			}

			httpx.OkJsonCtx(r.Context(), w, &types.ImageUploadRes{
				Base: types.Base{Code: 1, Msg: "文件上传成功"},
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

func Is(db *gorm.DB, hash string, path string) (is bool, info models.Upload) {
	if err := db.Model(&models.Upload{}).
		Select("id", "hash", "file_path", "origin_file_path").
		Where("hash = ? and origin_file_path = ?", hash, path).
		First(&info).Error; errors.As(err, &gorm.ErrRecordNotFound) {
		return false, info
	}
	return true, info
}
