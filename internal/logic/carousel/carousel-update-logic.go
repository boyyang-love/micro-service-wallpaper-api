package carousel

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CarouselUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCarouselUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CarouselUpdateLogic {
	return &CarouselUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CarouselUpdateLogic) CarouselUpdate(req *types.CarouselUpdateReq) (resp *types.CarouselUpdateRes, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.Carousel{}).
		Where("id = ?", req.Id).
		Select("path", "sort", "status").
		Updates(models.Carousel{
			Path:   req.Path,
			Sort:   req.Sort,
			Status: req.Status,
		}).
		Error; err != nil {
		return nil, err
	}

	return &types.CarouselUpdateRes{
		Base: types.Base{
			Code: 1,
			Msg:  "更新成功",
		},
	}, nil
}
