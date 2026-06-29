package album

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAlbumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAlbumLogic {
	return &UpdateAlbumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAlbumLogic) UpdateAlbum(req *types.AlbumUpdateReq) (resp *types.AlbumUpdateRes, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.Album{}).
		Select("name", "desc", "status", "cover").
		Where("id = ?", req.Id).
		Updates(models.Album{
			Name:   req.Name,
			Desc:   req.Desc,
			Status: req.Status,
			Cover:  req.Cover,
		}).Error; err != nil {
		return nil, err
	}

	return &types.AlbumUpdateRes{
		Base: types.Base{
			Code: 1,
			Msg:  "更新成功",
		},
	}, nil
}
