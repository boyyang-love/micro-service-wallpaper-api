package search

import (
	"context"
	"errors"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"gorm.io/gorm"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddSearchRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSearchRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSearchRecordsLogic {
	return &AddSearchRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSearchRecordsLogic) AddSearchRecords(req *types.AddSearchRecordsReq) (resp *types.AddSearchRecordsRes, err error) {
	var search models.Search
	if err = l.svcCtx.
		DB.
		Model(&models.Search{}).
		Where("keywords = ?", req.Keywords).
		First(&search).
		Error; err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			if err = l.svcCtx.
				DB.
				Model(&models.Search{}).
				Create(&models.Search{
					Keywords:    req.Keywords,
					SearchCount: 1,
				}).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if err = l.svcCtx.
		DB.
		Model(&models.Search{}).
		Where("keywords = ?", req.Keywords).
		Update("search_count", gorm.Expr("search_count + ?", 1)).
		Error; err != nil {
		return nil, err
	}

	return &types.AddSearchRecordsRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
	}, nil
}
