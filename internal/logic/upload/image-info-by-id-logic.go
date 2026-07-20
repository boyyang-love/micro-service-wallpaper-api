package upload

import (
	"context"
	"strings"

	"github.com/boyyang-love/micro-service-wallpaper-api/helper"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImageInfoByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImageInfoByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageInfoByIdLogic {
	return &ImageInfoByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImageInfoByIdLogic) ImageInfoById(req *types.ImageInfoByIdReq, authHeader string) (resp *types.ImageInfoByIdRes, err error) {
	var uploadInfo Upload
	var imageInfo types.ImageInfo
	if err = l.svcCtx.
		DB.
		Preload("Tags").
		Preload("Category").
		Preload("Recommend").
		Preload("Group").
		Model(&Upload{}).
		Where("id", req.Id).
		First(&uploadInfo).
		Error; err != nil {
		return nil, err
	}

	uploadInfo.Like = int(l.LikeNum(req.Id))

	_ = copier.Copy(&imageInfo, &uploadInfo)

	// 从 Authorization header 解析用户 ID（可选）
	userId := l.parseUserId(authHeader)
	if userId != "" {
		imageInfo.IsLiked = l.IsLiked(req.Id, userId)
	}

	return &types.ImageInfoByIdRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: imageInfo,
	}, nil
}

func (l *ImageInfoByIdLogic) parseUserId(authHeader string) string {
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := helper.ParseToken(token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return ""
	}
	return claims.Id
}

func (l *ImageInfoByIdLogic) LikeNum(id string) int64 {
	var count int64
	if err := l.svcCtx.
		DB.
		Model(&models.Like{}).
		Where("upload_id = ? and status = ?", id, 1).
		Count(&count).Error; err != nil {
		return 0
	}

	return count
}

func (l *ImageInfoByIdLogic) IsLiked(uploadId string, userId string) bool {
	var count int64
	if err := l.svcCtx.
		DB.
		Model(&models.Like{}).
		Where("upload_id = ? AND user_id = ? AND status = ?", uploadId, userId, true).
		Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}
