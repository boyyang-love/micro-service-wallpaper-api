package recommend

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecommendSortLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecommendSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecommendSortLogic {
	return &RecommendSortLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecommendSortLogic) RecommendSort(req *types.RecommendSortReq) (resp *types.RecommendSortRes, err error) {
	for _, sortData := range req.SortData {
		if err := l.svcCtx.
			DB.
			Model(&models.Recommend{}).
			Where("id=?", sortData.Id).
			Select("sort").
			Updates(&models.Recommend{Sort: sortData.Sort}).
			Error; err != nil {
			return nil, err
		}
	}

	return &types.RecommendSortRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
	}, nil

}
