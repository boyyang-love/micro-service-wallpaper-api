package category

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryInfoLogic {
	return &CategoryInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryInfoLogic) CategoryInfo(req *types.CategoryInfoReq) (resp *types.CategoryInfoRes, err error) {

	DB := l.svcCtx.DB.Model(&models.Category{}).Order("sort asc")

	if req.Name != "" {
		DB = DB.Where("name LIKE ?", "%"+req.Name+"%")
	}

	var categoryInfo []types.CategoryInfo
	var count int64
	if err = DB.
		Select("id", "name", "created", "updated", "sort", "web", "moa").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&categoryInfo).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	return &types.CategoryInfoRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.CategoryInfoData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: categoryInfo,
		},
	}, nil
}
