package upload

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImageInfoByHotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImageInfoByHotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageInfoByHotLogic {
	return &ImageInfoByHotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImageInfoByHotLogic) ImageInfoByHot(req *types.ImageInfoByHotReq) (resp *types.ImageInfoByHotRes, err error) {
	var count int64
	var info []types.ImageInfoHot
	if err := l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Order("view desc, download desc").
		Offset((req.Page-1)*req.Limit).
		Limit(req.Limit).
		Select("created", "updated", "id", "file_name", "file_path", "w", "h", "view").
		Where("status = ? and type = ?", 1, req.Type).
		Find(&info).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	return &types.ImageInfoByHotRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.ImageInfoByHotResdata{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: info,
		},
	}, nil
}
