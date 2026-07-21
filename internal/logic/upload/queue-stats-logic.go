package upload

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueueStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueueStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueueStatsLogic {
	return &QueueStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueueStatsLogic) QueueStats(req *types.QueueStatsReq) (resp *types.QueueStatsRes, err error) {
	return &types.QueueStatsRes{
		Base: types.Base{Code: 1, Msg: "ok"},
		Data: types.QueueStatsResData{
			WorkerCount: 0,
			QueueSize:   0,
			QueueLength: 0,
		},
	}, nil
}
