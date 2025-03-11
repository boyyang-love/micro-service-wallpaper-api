package category

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategorySummaryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategorySummaryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategorySummaryListLogic {
	return &CategorySummaryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type Category struct {
	models.Category
	Upload []Upload `json:"upload" gorm:"many2many:upload_category"`
}

type Upload struct {
	models.Upload
}

func (l *CategorySummaryListLogic) CategorySummaryList(req *types.CategorySummaryListReq) (resp *types.CategorySummaryListRes, err error) {
	var categories []Category
	var categorySummary []types.CategorySummary
	var count int64
	if err := l.svcCtx.
		DB.
		Debug().
		Order("sort").
		Model(Category{}).
		Select("id", "name", "sort").
		Preload("Upload", func(db *gorm.DB) *gorm.DB {
			return db.Select("id").Where("type = ?", req.Type)
		}).
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&categories).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	err = copier.Copy(&categorySummary, &categories)
	if err != nil {
		return nil, err
	}

	return &types.CategorySummaryListRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.CategorySummaryListData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: categorySummary,
		},
	}, nil
}
