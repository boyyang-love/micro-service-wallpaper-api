package category

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryUpdateLogic {
	return &CategoryUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryUpdateLogic) CategoryUpdate(req *types.CategoryUpdateReq) (resp *types.CategoryUpdateRes, err error) {

	if err = l.svcCtx.
		DB.
		Model(&models.Category{}).
		Where("id = ?", req.Id).
		Select("name", "sort", "web", "moa").
		Updates(&models.Category{
			Name: req.Name,
			Sort: req.Sort,
			Web:  req.Web,
			Moa:  req.Moa,
		}).Error; err != nil {
		return nil, err
	}

	return &types.CategoryUpdateRes{
		Base: types.Base{
			Code: 1,
			Msg:  "更新成功",
		},
	}, nil
}
