package category

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategorySortLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategorySortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategorySortLogic {
	return &CategorySortLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategorySortLogic) CategorySort(req *types.CategorySortReq) (resp *types.CategorySortRes, err error) {

	for _, sortData := range req.SortData {
		if err := l.svcCtx.
			DB.
			Model(&models.Category{}).
			Where("id=?", sortData.Id).
			Select("sort").
			Updates(&models.Category{Sort: sortData.Sort}).
			Error; err != nil {
			return nil, err
		}
	}

	return &types.CategorySortRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
	}, nil
}
