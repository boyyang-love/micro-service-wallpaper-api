package like

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeNumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeNumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeNumLogic {
	return &LikeNumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeNumLogic) LikeNum(req *types.LikeNumReq) (resp *types.LikeNumRes, err error) {
	var num int64
	if err = l.svcCtx.
		DB.
		Model(&models.Like{}).
		Where("upload_id = ? and status = ?", req.UploadId, true).
		Count(&num).
		Error; err != nil {
		return nil, err
	}

	return &types.LikeNumRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.LikeNumResData{Num: num},
	}, nil
}
