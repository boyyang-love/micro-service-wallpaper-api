package upload

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	uploadPb "github.com/boyyang-love/micro-service-wallpaper-rpc/upload/pb/upload"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChunkMergeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChunkMergeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChunkMergeLogic {
	return &ChunkMergeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChunkMergeLogic) ChunkMerge(req *types.ChunkMergeReq) (resp *types.ChunkMergeRes, err error) {
	logx.Infof("开始合并分片: 文件=%s, 哈希=%s, 分片数=%d", req.FileName, req.FileHash, req.TotalChunks)

	hashDir := filepath.Join(chunkDir, req.FileHash)

	tmpPath := filepath.Join(chunkDir, fmt.Sprintf("%s-merged-%d", req.FileHash, time.Now().Unix()))
	mergedFile, err := os.Create(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("创建合并文件失败: %v", err)
	}
	defer func() {
		mergedFile.Close()
		os.Remove(tmpPath)
	}()

	for i := 0; i < req.TotalChunks; i++ {
		chunkPath := filepath.Join(hashDir, fmt.Sprintf("%d", i))
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			return nil, fmt.Errorf("打开分片 %d 失败: %v", i, err)
		}

		if _, err := io.Copy(mergedFile, chunkFile); err != nil {
			chunkFile.Close()
			return nil, fmt.Errorf("合并分片 %d 失败: %v", i, err)
		}
		chunkFile.Close()
	}

	mergedFile.Seek(0, 0)

	hash := md5.New()
	if _, err := io.Copy(hash, mergedFile); err != nil {
		return nil, fmt.Errorf("计算哈希失败: %v", err)
	}
	mergedHash := hex.EncodeToString(hash.Sum(nil))

	mergedFile.Seek(0, 0)

	fileBytes, err := io.ReadAll(mergedFile)
	if err != nil {
		return nil, fmt.Errorf("读取合并文件失败: %v", err)
	}

	logx.Infof("分片合并成功: 文件=%s, 大小=%d, 哈希=%s", req.FileName, len(fileBytes), mergedHash)

	ext := filepath.Ext(req.FileName)
	oriPath := fmt.Sprintf("%s/%s/%s/original/%s%s", req.RootDir, req.Dir, req.UserId, mergedHash, ext)
	comPath := fmt.Sprintf("%s/%s/%s/compress/%s.webp", req.RootDir, req.Dir, req.UserId, mergedHash)

	imageUpload, err := l.svcCtx.UploadService.CosUpload(
		l.ctx,
		&uploadPb.ImageUploadReq{
			File:       fileBytes,
			Path:       comPath,
			OriPath:    oriPath,
			Quality:    req.Quality,
			BucketName: req.BucketName,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("上传到云存储失败: %v", err)
	}

	uploadInfo := models.Upload{
		Hash:           imageUpload.Data.OriETag,
		FileName:       req.FileName,
		OriginFileSize: int64(imageUpload.Data.OriSize),
		FileSize:       int64(imageUpload.Data.Size),
		OriginType:     strings.TrimPrefix(ext, "."),
		FileType:       "webp",
		OriginFilePath: oriPath,
		FilePath:       comPath,
		Type:           req.Type,
		Status:         req.Status,
		UserId:         req.UserId,
	}
	if err := l.svcCtx.DB.Model(&models.Upload{}).Create(&uploadInfo).Error; err != nil {
		return nil, fmt.Errorf("保存到数据库失败: %v", err)
	}

	l.saveRelations(uploadInfo.Id, req)

	go func() {
		os.RemoveAll(hashDir)
	}()

	return &types.ChunkMergeRes{
		Base: types.Base{Code: 1, Msg: "分片合并成功"},
		Data: types.ChunkMergeResData{
			FileHash:    mergedHash,
			FileName:    req.FileName,
			FileSize:    len(fileBytes),
			TotalChunks: req.TotalChunks,
		},
	}, nil
}

func (l *ChunkMergeLogic) saveRelations(uploadId string, req *types.ChunkMergeReq) {
	if req.Tags != "" {
		var uploadTags []models.UploadTag
		for _, v := range strings.Split(req.Tags, ",") {
			uploadTags = append(uploadTags, models.UploadTag{UploadId: uploadId, TagId: v})
		}
		l.svcCtx.DB.Model(&models.UploadTag{}).Create(&uploadTags)
	}

	if req.Category != "" {
		var uploadCategory []models.UploadCategory
		for _, v := range strings.Split(req.Category, ",") {
			uploadCategory = append(uploadCategory, models.UploadCategory{UploadId: uploadId, CategoryId: v})
		}
		l.svcCtx.DB.Model(&models.UploadCategory{}).Create(&uploadCategory)
	}

	if req.Recommend != "" {
		var uploadRecommend []models.UploadRecommend
		for _, v := range strings.Split(req.Recommend, ",") {
			uploadRecommend = append(uploadRecommend, models.UploadRecommend{UploadId: uploadId, RecommendId: v})
		}
		l.svcCtx.DB.Model(&models.UploadRecommend{}).Create(&uploadRecommend)
	}

	if req.Group != "" {
		var uploadGroup []models.UploadGroup
		for _, v := range strings.Split(req.Group, ",") {
			uploadGroup = append(uploadGroup, models.UploadGroup{UploadId: uploadId, GroupId: v})
		}
		l.svcCtx.DB.Model(&models.UploadGroup{}).Create(&uploadGroup)
	}

	if req.Album != "" {
		var uploadAlbum []models.UploadAlbum
		for _, v := range strings.Split(req.Album, ",") {
			uploadAlbum = append(uploadAlbum, models.UploadAlbum{UploadId: uploadId, AlbumId: v})
		}
		l.svcCtx.DB.Model(&models.UploadAlbum{}).Create(&uploadAlbum)
	}
}
