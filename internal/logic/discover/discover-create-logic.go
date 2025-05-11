package discover

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DiscoverCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverCreateLogic {
	return &DiscoverCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DiscoverCreateLogic) DiscoverCreate(req *types.DiscoverCreateReq) (resp *types.DiscoverCreateRes, err error) {
	userid := fmt.Sprintf("%s", l.ctx.Value("Id"))
	if err := l.svcCtx.DB.Model(&models.Discover{}).Create(&models.Discover{
		Title:    req.Title,
		Subtitle: req.Subtitle,
		UserId:   userid,
		ImageIds: req.ImageIds,
	}).Error; err != nil {
		return nil, err
	}

	return &types.DiscoverCreateRes{
		Base: types.Base{
			Code: 1,
			Msg:  "上传成功",
		},
	}, nil
}
