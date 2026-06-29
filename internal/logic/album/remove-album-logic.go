package album

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveAlbumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveAlbumLogic {
	return &RemoveAlbumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveAlbumLogic) RemoveAlbum(req *types.AlbumRemoveReq) (resp *types.AlbumRemoveRes, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.Album{}).
		Delete(&models.Album{}).
		Where("id = ?", req.Id).
		Error; err != nil {
		return nil, err
	}

	return &types.AlbumRemoveRes{
		Base: types.Base{
			Code: 1,
			Msg:  "删除成功",
		},
	}, nil
}
