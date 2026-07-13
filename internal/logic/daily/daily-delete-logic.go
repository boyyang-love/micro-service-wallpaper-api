package daily

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type DailyDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDailyDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DailyDeleteLogic {
	return &DailyDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DailyDeleteLogic) DailyDelete(req *types.DailyDeleteReq) (resp *types.DailyDeleteRes, err error) {
	if err := l.svcCtx.DB.
		Where("id = ?", req.Id).
		Delete(&models.DailyWallpaper{}).Error; err != nil {
		return nil, err
	}

	return &types.DailyDeleteRes{
		Base: types.Base{Code: 1, Msg: "ok"},
	}, nil
}
