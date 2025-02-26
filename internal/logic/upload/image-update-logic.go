package upload

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImageUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImageUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageUpdateLogic {
	return &ImageUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImageUpdateLogic) ImageUpdate(req *types.ImageUpdateReq) (resp *types.ImageUpdateRes, err error) {

	if err := l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Select("file_name", "type", "status").
		Where("id = ?", req.Id).
		Updates(&models.Upload{
			FileName: req.FileName,
			Type:     req.Type,
			Status:   req.Status,
		}).
		Error; err != nil {
		return nil, err
	}
	return &types.ImageUpdateRes{
		Base: types.Base{
			Code: 1,
			Msg:  "更新成功",
		},
	}, nil
}
