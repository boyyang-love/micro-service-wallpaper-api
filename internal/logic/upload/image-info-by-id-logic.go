package upload

import (
	"context"
	"github.com/jinzhu/copier"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

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

func (l *ImageInfoByIdLogic) ImageInfoById(req *types.ImageInfoByIdReq) (resp *types.ImageInfoByIdRes, err error) {
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

	_ = copier.Copy(&imageInfo, &uploadInfo)

	return &types.ImageInfoByIdRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: imageInfo,
	}, nil
}
