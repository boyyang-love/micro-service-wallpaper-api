package group

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGroupLogic {
	return &ListGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListGroupLogic) ListGroup(req *types.GroupListReq) (resp *types.GroupListRes, err error) {
	var records []types.GroupListInfo
	var count int64

	DB := l.svcCtx.DB
	if req.Name != "" {
		DB = DB.Where("name like ?", "%"+req.Name+"%")
	}

	if err = DB.
		Order("created desc").
		Model(&models.Group{}).
		Select("created", "updated", "id", "name").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&records).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	return &types.GroupListRes{
		Base: types.Base{
			Code: 1,
			Msg:  "OK",
		},
		Data: types.GroupListData{
			BaseRecord: types.BaseRecord{
				Total: count,
				Limit: req.Limit,
				Page:  req.Page,
			},
			Records: records,
		},
	}, nil
}
