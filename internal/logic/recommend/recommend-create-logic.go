package recommend

import (
	"context"
	"errors"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"gorm.io/gorm"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecommendCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecommendCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecommendCreateLogic {
	return &RecommendCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecommendCreateLogic) RecommendCreate(req *types.RecommendCreateReq) (resp *types.RecommendCreateRes, err error) {

	is, err := l.IsExist(req)
	if err != nil {
		return nil, err
	}

	if is {
		return &types.RecommendCreateRes{
			Base: types.Base{
				Code: 2,
				Msg:  "该名称已经存在",
			},
		}, nil
	}

	if err = l.svcCtx.
		DB.
		Model(&models.Recommend{}).
		Create(&models.Recommend{
			Name: req.Name,
			Sort: req.Sort,
		}).
		Error; err != nil {
		return nil, err
	}

	return &types.RecommendCreateRes{
		Base: types.Base{
			Code: 1,
			Msg:  "新增成功",
		},
	}, nil
}

func (l *RecommendCreateLogic) IsExist(req *types.RecommendCreateReq) (is bool, err error) {
	var recommend models.Recommend
	if err := l.svcCtx.
		DB.
		Model(&models.Recommend{}).
		Select("name").
		Where("name = ?", req.Name).
		First(&recommend).
		Error; err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			return false, nil
		}
		return true, err
	}

	return true, nil
}
