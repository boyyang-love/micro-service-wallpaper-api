package block

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveBlockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveBlockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveBlockLogic {
	return &RemoveBlockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveBlockLogic) RemoveBlock(req *types.BlockRemoveReq) (resp *types.BlockRemoveRes, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.Block{}).
		Where("id = ?", req.Id).
		Delete(&models.Block{}).
		Error; err != nil {
		return nil, err
	}

	return &types.BlockRemoveRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
	}, nil
}
