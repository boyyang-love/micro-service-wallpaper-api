package recommend

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecommendRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecommendRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecommendRemoveLogic {
	return &RecommendRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecommendRemoveLogic) RecommendRemove(req *types.RecommendRemoveReq) (resp *types.RecommendRemoveRes, err error) {

	if err = l.svcCtx.
		DB.
		Model(&models.Recommend{}).
		Where("id = ?", req.Id).
		Delete(&models.Recommend{}).
		Error; err != nil {
		return nil, err
	}

	return &types.RecommendRemoveRes{
		Base: types.Base{
			Code: 1,
			Msg:  "删除成功",
		},
	}, nil
}
