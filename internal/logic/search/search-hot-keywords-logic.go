package search

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/jinzhu/copier"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchHotKeywordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchHotKeywordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchHotKeywordsLogic {
	return &SearchHotKeywordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchHotKeywordsLogic) SearchHotKeywords(req *types.SearchHotKeywordsReq) (resp *types.SearchHotKeywordsRes, err error) {
	search := make([]models.Search, 0)
	records := make([]types.SearchHotKeywordsInfo, 0)
	var count int64
	if err = l.svcCtx.
		DB.
		Order("search_count desc").
		Model(&models.Search{}).
		Select("id", "created", "updated", "keywords", "search_count").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&search).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	err = copier.Copy(&records, &search)
	if err != nil {
		return nil, err
	}

	return &types.SearchHotKeywordsRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.SearchHotKeywordsData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: records,
		},
	}, nil
}
