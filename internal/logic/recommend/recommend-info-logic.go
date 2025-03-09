package recommend

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecommendInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecommendInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecommendInfoLogic {
	return &RecommendInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecommendInfoLogic) RecommendInfo(req *types.RecommendInfoReq) (resp *types.RecommendInfoRes, err error) {
	DB := l.svcCtx.
		DB.
		Order("sort asc").
		Model(&models.Recommend{})

	if req.Name != "" {
		DB = DB.Where("name LIKE ?", "%"+req.Name+"%")
	}

	var recommendInfo []types.RecommendInfo
	var count int64
	if err = DB.
		Select("id", "name", "created", "updated", "sort").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&recommendInfo).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	return &types.RecommendInfoRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.RecommendInfoData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: recommendInfo,
		},
	}, nil
}
