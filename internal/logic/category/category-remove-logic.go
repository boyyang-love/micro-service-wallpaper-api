package category

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryRemoveLogic {
	return &CategoryRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryRemoveLogic) CategoryRemove(req *types.CategoryRemoveReq) (resp *types.CategoryRemoveRes, err error) {

	if err = l.svcCtx.
		DB.
		Model(&models.Category{}).
		Where("id = ?", req.Id).
		Delete(&models.Category{}).
		Error; err != nil {
		return nil, err
	}
	return &types.CategoryRemoveRes{Base: types.Base{
		Code: 1,
		Msg:  "删除成功",
	}}, nil
}
