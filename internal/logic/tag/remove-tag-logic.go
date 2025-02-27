package tag

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveTagLogic {
	return &RemoveTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveTagLogic) RemoveTag(req *types.RemoveTagReq) (resp *types.RemoveTagRes, err error) {
	if err := l.svcCtx.
		DB.
		Model(&models.Tag{}).
		Where("id = ?", req.Id).
		Delete(&models.Tag{}).
		Error; err != nil {
		return nil, err
	}

	return &types.RemoveTagRes{
		Base: types.Base{
			Code: 1,
			Msg:  "删除成功",
		},
	}, nil
}
