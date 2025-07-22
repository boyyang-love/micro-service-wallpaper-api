package like

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeRemoveLogic {
	return &LikeRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeRemoveLogic) LikeRemove(req *types.LikeRemoveReq) (resp *types.LikeRemoveRes, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.Like{}).
		Where("id = ?", req.Id).
		Delete(&models.Like{}).
		Error; err != nil {
		return nil, err
	}

	return &types.LikeRemoveRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
	}, nil
}
