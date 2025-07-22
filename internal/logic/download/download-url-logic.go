package download

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadUrlLogic {
	return &DownloadUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadUrlLogic) DownloadUrl(req *types.DownloadUrlReq) (resp *types.DownloadUrlRes, err error) {
	var data types.DownloadUrlData
	if err := l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Select("id", "origin_file_path", "file_name", "origin_type", "type").
		Where("id = ?", req.Id).
		First(&data).
		Error; err != nil {
		return nil, err
	}

	if err = l.AddDownloadRecord(req, data.Type); err != nil {
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

func (l *DownloadUrlLogic) AddDownloadRecord(req *types.DownloadUrlReq, downloadType string) error {
	userid := fmt.Sprintf("%s", l.ctx.Value("Id"))

	if err := l.svcCtx.
		DB.
		Model(&models.Download{}).
		Create(&models.Download{
			DownloadId: req.Id,
			UserId:     userid,
			Type:       downloadType,
		}).Error; err != nil {
		return err
	}

	return nil
}
