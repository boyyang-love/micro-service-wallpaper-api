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

type RecommendUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecommendUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecommendUpdateLogic {
	return &RecommendUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecommendUpdateLogic) RecommendUpdate(req *types.RecommendUpdateReq) (resp *types.RecommendUpdateRes, err error) {
	is, err := l.IsExist(req)
	if err != nil {
		return nil, err
	}

	if is {
		return &types.RecommendUpdateRes{Base: types.Base{
			Code: 2,
			Msg:  "名称已存在",
		}}, nil
	}

	if err = l.svcCtx.
		DB.
		Model(&models.Recommend{}).
		Where("id = ?", req.Id).
		Updates(&models.Recommend{
			Name: req.Name,
			Sort: req.Sort,
		}).
		Error; err != nil {
		return nil, err
	}
	return &types.RecommendUpdateRes{Base: types.Base{
		Code: 1,
		Msg:  "修改成功",
	}}, nil
}

func (l *RecommendUpdateLogic) IsExist(req *types.RecommendUpdateReq) (is bool, err error) {
	var recommend models.Recommend
	if err := l.svcCtx.
		DB.
		Model(&models.Recommend{}).
		Select("name").
		Where("name = ? and id != ?", req.Name, req.Id).
		First(&recommend).
		Error; err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			return false, nil
		}
		return true, err
	}

	return true, nil
}
