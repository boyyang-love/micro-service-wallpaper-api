package carousel

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CarouselCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCarouselCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CarouselCreateLogic {
	return &CarouselCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CarouselCreateLogic) CarouselCreate(req *types.CarouselCreateReq) (resp *types.CarouselCreateRes, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.Carousel{}).
		Create(&models.Carousel{
			Path:   req.Path,
			Sort:   req.Sort,
			Status: req.Status,
		}).Error; err != nil {
		return nil, err
	}

	return &types.CarouselCreateRes{
		Base: types.Base{
			Code: 1,
			Msg:  "创建成功",
		},
	}, nil
}
