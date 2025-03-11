package upload

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"gorm.io/gorm"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImageSummaryUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImageSummaryUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageSummaryUpdateLogic {
	return &ImageSummaryUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImageSummaryUpdateLogic) ImageSummaryUpdate(req *types.ImageSummaryUpdateReq) (resp *types.ImageSummaryUpdateRes, err error) {

	DB := l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Where("id = ?", req.Id)

	if req.Type == "download" {
		DB = DB.Update("download", gorm.Expr("download + ?", 1))
	}

	if req.Type == "view" {
		DB = DB.Update("view", gorm.Expr("view + ?", 1))
	}

	if err := DB.Error; err != nil {
		return nil, err
	}

	return &types.ImageSummaryUpdateRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
	}, nil
}
