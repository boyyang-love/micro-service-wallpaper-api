// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package like

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeRemovePostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeRemovePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeRemovePostLogic {
	return &LikeRemovePostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeRemovePostLogic) LikeRemovePost(req *types.PostLikeRemoveReq) (resp *types.PostLikeRemoveRes, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.Like{}).
		Where("upload_id = ?", req.UploadId).
		Delete(&models.Like{}).
		Error; err != nil {
		return nil, err
	}

	return &types.PostLikeRemoveRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
	}, nil
}
