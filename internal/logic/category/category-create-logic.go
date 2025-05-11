package category

import (
	"context"
	"errors"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"gorm.io/gorm"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryCreateLogic {
	return &CategoryCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryCreateLogic) CategoryCreate(req *types.CategoryCreateReq) (resp *types.CategoryCreateRes, err error) {

	is, err := l.IsExist(req)
	if err != nil {
		return nil, err
	}

	if is {
		return &types.CategoryCreateRes{
			Base: types.Base{
				Code: 2,
				Msg:  "该分类已经存在",
			},
		}, nil
	}

	if err = l.svcCtx.
		DB.
		Model(&models.Category{}).
		Create(&models.Category{
			Name: req.Name,
			Sort: req.Sort,
			Web:  req.Web,
			Moa:  req.Moa,
		}).
		Error; err != nil {
		return nil, err
	}

	return &types.CategoryCreateRes{
		Base: types.Base{
			Code: 1,
			Msg:  "创建成功",
		},
	}, nil
}

func (l *CategoryCreateLogic) IsExist(req *types.CategoryCreateReq) (is bool, err error) {
	var category models.Category
	if err := l.svcCtx.
		DB.
		Model(&models.Category{}).
		Select("name").
		Where("name = ?", req.Name).
		First(&category).
		Error; err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			return false, nil
		}
		return true, err
	}

	return true, nil
}
