package upload

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImageSummaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImageSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageSummaryLogic {
	return &ImageSummaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImageSummaryLogic) ImageSummary(req *types.ImageSummaryReq) (resp *types.ImageSummaryRes, err error) {
	var pc int64
	var moa int64
	if err := l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Select("id").
		Where("status = ? and type = ?", 1, "PC").
		Count(&pc).
		Error; err != nil {
		return nil, err
	}

	if err := l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Select("id").
		Where("status = ? and type = ?", 1, "MOA").
		Count(&moa).
		Error; err != nil {
		return nil, err
	}

	return &types.ImageSummaryRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.ImageSummaryResdata{
			Pc:  pc,
			Moa: moa,
		},
	}, nil
}
