package carousel

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CarouselListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCarouselListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CarouselListLogic {
	return &CarouselListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CarouselListLogic) CarouselList(req *types.CarouselListReq) (resp *types.CarouselListRes, err error) {
	var list []types.CarouselInfo
	var count int64

	DB := l.svcCtx.
		DB.
		Model(&models.Carousel{}).
		Order("sort asc")

	if req.Status != 0 {
		DB = DB.Where("status = ?", req.Status)
	}

	if err := DB.
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&list).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	return &types.CarouselListRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.CarouselListData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: list,
		},
	}, nil
}
