package upload

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"

	"github.com/boyyang-love/micro-service-wallpaper-api/helper"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	upload2 "github.com/boyyang-love/micro-service-wallpaper-rpc/upload/pb/upload"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ImageUploadAsyncHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		logx.Infof("用户 %s 异步上传文件: %s, 大小: %d bytes", userid, fileHeader.Filename, fileHeader.Size)

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
		}

		// 生成任务ID
		taskId := uuid.New().String()

		// 保存必要的数据副本用于异步处理
		fileBytes := make([]byte, originBuffer.Len())
		copy(fileBytes, originBuffer.Bytes())
		uploadReq := req
		userId := userid
		imgWidth := img.Bounds().Dx()
		imgHeight := img.Bounds().Dy()
		svcCtxRef := svcCtx

		// 立即返回任务ID
		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
			"code": 1,
			"msg":  "上传任务已提交，请稍后查询结果",
			"data": map[string]interface{}{
				"taskId":   taskId,
				"fileName": req.FileName,
				"status":   "processing",
			},
		})

		// 异步处理上传（使用独立的context）
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logx.Errorf("异步上传任务 %s 发生panic: %v", taskId, r)
				}
			}()

			// 使用Background context，不依赖HTTP请求的生命周期
			ctx := context.Background()

			logx.Infof("开始异步处理上传任务: %s", taskId)

			// 上传到云存储
			imageUpload, err := svcCtxRef.UploadService.CosUpload(
				ctx,
				&upload2.ImageUploadReq{
					File:       fileBytes,
					Path:       comPath,
					OriPath:    oriPath,
					Quality:    uploadReq.Quality,
					BucketName: uploadReq.BucketName,
				},
			)
			if err != nil {
				logx.Errorf("异步上传任务 %s 上传到云存储失败: %v", taskId, err)
				return
			}

			// 保存到数据库
			uploadInfo := models.Upload{
				Hash:           imageUpload.Data.OriETag,
				FileName:       uploadReq.FileName,
				OriginFileSize: int64(imageUpload.Data.OriSize),
				FileSize:       int64(imageUpload.Data.Size),
				OriginType:     imgType,
				FileType:       "webp",
				OriginFilePath: oriPath,
				FilePath:       comPath,
				Type:           uploadReq.Type,
				W:              imgWidth,
				H:              imgHeight,
				Status:         uploadReq.Status,
				UserId:         userId,
			}
			if err := svcCtxRef.DB.Model(&models.Upload{}).Create(&uploadInfo).Error; err != nil {
				logx.Errorf("异步上传任务 %s 保存到数据库失败: %v", taskId, err)
				return
			}

			// 保存关联数据
			saveAsyncUploadRelations(svcCtxRef, uploadInfo.Id, uploadReq)

			logx.Infof("异步上传任务 %s 完成, 文件ID: %s", taskId, uploadInfo.Id)
		}()
	}
}

func saveAsyncUploadRelations(svcCtx *svc.ServiceContext, uploadId string, req types.ImageUploadReq) {
	if req.Tags != "" {
		var uploadTags []models.UploadTag
		for _, v := range strings.Split(req.Tags, ",") {
			uploadTags = append(uploadTags, models.UploadTag{UploadId: uploadId, TagId: v})
		}
		svcCtx.DB.Model(&models.UploadTag{}).Create(&uploadTags)
	}

	if req.Category != "" {
		var uploadCategory []models.UploadCategory
		for _, v := range strings.Split(req.Category, ",") {
			uploadCategory = append(uploadCategory, models.UploadCategory{UploadId: uploadId, CategoryId: v})
		}
		svcCtx.DB.Model(&models.UploadCategory{}).Create(&uploadCategory)
	}

	if req.Recommend != "" {
		var uploadRecommend []models.UploadRecommend
		for _, v := range strings.Split(req.Recommend, ",") {
			uploadRecommend = append(uploadRecommend, models.UploadRecommend{UploadId: uploadId, RecommendId: v})
		}
		svcCtx.DB.Model(&models.UploadRecommend{}).Create(&uploadRecommend)
	}

	if req.Group != "" {
		var uploadGroup []models.UploadGroup
		for _, v := range strings.Split(req.Group, ",") {
			uploadGroup = append(uploadGroup, models.UploadGroup{UploadId: uploadId, GroupId: v})
		}
		svcCtx.DB.Model(&models.UploadGroup{}).Create(&uploadGroup)
	}

	if req.Album != "" {
		var uploadAlbum []models.UploadAlbum
		for _, v := range strings.Split(req.Album, ",") {
			uploadAlbum = append(uploadAlbum, models.UploadAlbum{UploadId: uploadId, AlbumId: v})
		}
		svcCtx.DB.Model(&models.UploadAlbum{}).Create(&uploadAlbum)
	}
}
