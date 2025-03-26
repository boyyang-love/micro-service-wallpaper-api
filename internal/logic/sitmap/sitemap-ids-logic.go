package sitmap

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SitemapIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSitemapIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SitemapIdsLogic {
	return &SitemapIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SitemapIdsLogic) SitemapIds(req *types.SitemapReq) (resp *types.SitemapRes, err error) {
	var ids []string
	var count int64
	if err = l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Select("id").
		Where("type = ?", req.Type).
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&ids).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	return &types.SitemapRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.SitemapResData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: ids,
		},
	}, nil
}
