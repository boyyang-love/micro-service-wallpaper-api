package daily

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type DailyUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDailyUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DailyUpdateLogic {
	return &DailyUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DailyUpdateLogic) DailyUpdate(req *types.DailyUpdateReq) (resp *types.DailyUpdateRes, err error) {
	updates := map[string]interface{}{}
	if req.UploadId != "" {
		updates["upload_id"] = req.UploadId
	}
	if req.Date != "" {
		updates["date"] = req.Date
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	updates["status"] = req.Status

	if err := l.svcCtx.DB.
		Model(&models.DailyWallpaper{}).
		Where("id = ?", req.Id).
		Updates(updates).Error; err != nil {
		return nil, err
	}

	return &types.DailyUpdateRes{
		Base: types.Base{Code: 1, Msg: "ok"},
	}, nil
}
