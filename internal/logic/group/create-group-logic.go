package group

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateGroupLogic) CreateGroup(req *types.GroupCreateReq) (resp *types.GroupCreateRes, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.Group{}).
		Where("name = ?", req.Name).
		FirstOrCreate(&models.Group{
			Name: req.Name,
		}).
		Error; err != nil {
		return nil, err
	}

	return &types.GroupCreateRes{
		Base: types.Base{
			Code: 1,
			Msg:  "创建成功",
		},
	}, nil
}
