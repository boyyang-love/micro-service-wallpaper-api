package upload

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImageInfoByCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImageInfoByCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageInfoByCategoryLogic {
	return &ImageInfoByCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImageInfoByCategoryLogic) ImageInfoByCategory(req *types.ImageInfoByCategoryReq) (resp *types.ImageInfoByCategoryRes, err error) {

	DB := l.svcCtx.
		DB.
		Order("created desc").
		Model(&models.Upload{}).
		Select("id", "file_name", "file_path", "w", "h", "created", "updated")

	if req.CategoryId != "" {
		ids, err := l.UploadIds(req)
		if err != nil {
			return nil, err
		}

		if len(*ids) == 0 {
			return &types.ImageInfoByCategoryRes{
				Base: types.Base{
					Code: 1,
					Msg:  "ok",
				},
				Data: types.ImageInfoByCategoryResdata{
					BaseRecord: types.BaseRecord{
						Page:  req.Page,
						Limit: req.Limit,
						Total: 0,
					},
					Records: []types.ImageInfoCategory{},
				},
			}, nil
		}

		DB = DB.Where("id IN (?)", *ids)
	}

	var uploadInfo []types.ImageInfoCategory
	var count int64
	if err := DB.
		Where("type = ? and status = ?", req.Type, 1).
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&uploadInfo).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}
	return &types.ImageInfoByCategoryRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.ImageInfoByCategoryResdata{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: uploadInfo,
		},
	}, nil
}

func (l *ImageInfoByCategoryLogic) UploadIds(req *types.ImageInfoByCategoryReq) (ids *[]string, err error) {
	var uploadIds []string
	if err = l.svcCtx.
		DB.
		Model(&models.UploadCategory{}).
		Select("upload_id").
		Where("category_id = ?", req.CategoryId).
		Find(&uploadIds).
		Error; err != nil {
		return nil, err
	}

	return &uploadIds, nil
}
