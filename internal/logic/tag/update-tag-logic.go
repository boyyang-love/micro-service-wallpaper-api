package tag

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTagLogic {
	return &UpdateTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTagLogic) UpdateTag(req *types.UpdateTagReq) (resp *types.UpdateTagRes, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.Tag{}).
		Where("id = ?", req.Id).
		Select("name", "type", "sort").
		Updates(&models.Tag{Name: req.Name, Type: req.Type, Sort: req.Sort}).
		Error; err != nil {
		return nil, err
	}

	return &types.UpdateTagRes{
		Base: types.Base{
			Code: 1,
			Msg:  "更新成功",
		},
	}, nil
}
