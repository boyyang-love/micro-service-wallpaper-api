package like

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeCreateLogic {
	return &LikeCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeCreateLogic) LikeCreate(req *types.LikeCreateOrUpdateReq) (resp *types.LikeCreateOrUpdateRes, err error) {
	// todo: add your logic here and delete this line

	return
}
