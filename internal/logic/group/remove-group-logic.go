package group

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveGroupLogic {
	return &RemoveGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveGroupLogic) RemoveGroup(req *types.GroupRemoveReq) (resp *types.GroupRemoveRes, err error) {

	if err = l.svcCtx.
		DB.
		Model(&models.Group{}).
		Where("id = ?", req.Id).
		Delete(&models.Group{}).
		Error; err != nil {
		return nil, err
	}

	return &types.GroupRemoveRes{
		Base: types.Base{
			Code: 1,
			Msg:  "删除成功",
		},
	}, nil
}
