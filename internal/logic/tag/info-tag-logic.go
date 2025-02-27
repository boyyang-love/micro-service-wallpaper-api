package tag

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoTagLogic {
	return &InfoTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoTagLogic) InfoTag(req *types.TagInfoReq) (resp *types.TagInfoRes, err error) {
	var records []types.TagInfo
	var total int64

	DB := l.svcCtx.
		DB.
		Order("created desc").
		Model(&models.Tag{}).
		Select("id", "name", "type", "created", "updated")

	if req.Name != "" {
		DB = DB.Where("name LIKE ?", "%"+req.Name+"%")
	}

	if req.Type != "" {
		DB = DB.Where("type = ?", req.Type)
	}

	if err = DB.
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Find(&records).
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	return &types.TagInfoRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.TagInfoResData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: total,
			},
			Records: records,
		},
	}, nil
}
