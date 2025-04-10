package download

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadUrlNoTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadUrlNoTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadUrlNoTokenLogic {
	return &DownloadUrlNoTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadUrlNoTokenLogic) DownloadUrlNoToken(req *types.DownloadUrlReq) (resp *types.DownloadUrlRes, err error) {
	var data types.DownloadUrlData
	if err := l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Select("origin_file_path", "file_name", "origin_type").
		Where("id = ?", req.Id).
		First(&data).
		Error; err != nil {
		return nil, err
	}

	return &types.DownloadUrlRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: data,
	}, nil
}
