package discover

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DiscoverUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverUpdateLogic {
	return &DiscoverUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DiscoverUpdateLogic) DiscoverUpdate(req *types.DiscoverUpdateStatusReq) (resp *types.DiscoverUpdateStatusRes, err error) {

	if err := l.svcCtx.
		DB.
		Model(&models.Discover{}).
		Where("id = ?", req.Id).
		Update("status", req.Status).
		Error; err != nil {
		return nil, err
	}

	return &types.DiscoverUpdateStatusRes{
		Base: types.Base{
			Code: 1,
			Msg:  "更新成功",
		},
	}, nil
}
