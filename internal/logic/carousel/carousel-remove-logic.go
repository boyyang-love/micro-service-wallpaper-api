package carousel

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CarouselRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCarouselRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CarouselRemoveLogic {
	return &CarouselRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CarouselRemoveLogic) CarouselRemove(req *types.CarouselRemoveReq) (resp *types.CarouselRemoveRes, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.Carousel{}).
		Where("id = ?", req.Id).
		Delete(&models.Carousel{}).
		Error; err != nil {
		return nil, err
	}

	return &types.CarouselRemoveRes{
		Base: types.Base{
			Code: 1,
			Msg:  "删除成功",
		},
	}, nil
}
