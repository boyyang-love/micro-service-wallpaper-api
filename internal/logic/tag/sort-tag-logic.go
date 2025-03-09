package tag

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SortTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSortTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SortTagLogic {
	return &SortTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SortTagLogic) SortTag(req *types.TagSortReq) (resp *types.TagSortRes, err error) {
	for _, sortData := range req.SortData {
		if err := l.svcCtx.
			DB.
			Model(&models.Tag{}).
			Where("id=?", sortData.Id).
			Select("sort").
			Updates(&models.Tag{Sort: sortData.Sort}).
			Error; err != nil {
			return nil, err
		}
	}

	return &types.TagSortRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
	}, nil
}
