package upload

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/zeromicro/go-zero/core/logx"
)

type ImageInfoByGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImageInfoByGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageInfoByGroupLogic {
	return &ImageInfoByGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImageInfoByGroupLogic) ImageInfoByGroup(req *types.ImageInfoByGroupReq) (resp *types.ImageInfoByGroupRes, err error) {
	var records []types.ImageInfoGroup
	ids, err := l.UploadIds(req)

	if err != nil {
		return nil, err
	}

	if err = l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Select("created", "updated", "id", "file_name", "file_path", "w", "h").
		Where("id in (?) and type = ?", ids, req.Type).
		Find(&records).
		Error; err != nil {
		return nil, err
	}

	return &types.ImageInfoByGroupRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.ImageInfoByGroupResdata{
			Records: records,
		},
	}, nil
}

func (l *ImageInfoByGroupLogic) UploadIds(req *types.ImageInfoByGroupReq) (ids []string, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.UploadGroup{}).
		Select("upload_id").
		Where("group_id = ?", req.GroupId).
		Find(&ids).
		Error; err != nil {
		return nil, err
	}

	return ids, nil
}
