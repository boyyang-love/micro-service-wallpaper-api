package album

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAlbumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAlbumLogic {
	return &CreateAlbumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAlbumLogic) CreateAlbum(req *types.AlbumCreateReq) (resp *types.AlbumCreateRes, err error) {

	if err = l.svcCtx.
		DB.
		Model(&models.Album{}).
		Create(&models.Album{
			Name:   req.Name,
			Desc:   req.Desc,
			Cover:  req.Cover,
			Status: req.Status,
		}).Error; err != nil {
		return nil, err
	}
	return &types.AlbumCreateRes{
		Base: types.Base{
			Code: 1,
			Msg:  "创建成功",
		},
	}, nil
}
