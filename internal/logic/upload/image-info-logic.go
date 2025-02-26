package upload

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImageInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImageInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageInfoLogic {
	return &ImageInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImageInfoLogic) ImageInfo(req *types.ImageInfoReq) (resp *types.ImageInfoRes, err error) {
	var uploadInfo []types.ImageInfo
	var count int64
	DB := l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Order("created  desc")

	if req.Status == 1 || req.Status == 2 {
		DB = DB.Where("status = ? ", req.Status)
	}

	if req.FileName != "" {
		DB = DB.Where("file_name LIKE ? ", "%"+req.FileName+"%")
	}
	if req.Type != "" {
		DB = DB.Where("type = ? ", req.Type)
	}

	if err := DB.
		Count(&count).
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Find(&uploadInfo).
		Offset(-1).
		Error; err != nil {
		return nil, err
	}

	return &types.ImageInfoRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.ImageInfoResdata{
			Page:    req.Page,
			Limit:   req.Limit,
			Total:   count,
			Records: uploadInfo,
		},
	}, nil
}
