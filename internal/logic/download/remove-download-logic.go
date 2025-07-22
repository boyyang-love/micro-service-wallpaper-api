package download

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveDownloadLogic {
	return &RemoveDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveDownloadLogic) RemoveDownload(req *types.DownloadRemoveReq) (resp *types.DownloadRemoveRes, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.Download{}).
		Where("id = ?", req.Id).
		Delete(&models.Download{}).
		Error; err != nil {
		return nil, err
	}

	return &types.DownloadRemoveRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
	}, nil
}
