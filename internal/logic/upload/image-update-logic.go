package upload

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImageUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImageUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageUpdateLogic {
	return &ImageUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImageUpdateLogic) ImageUpdate(req *types.ImageUpdateReq) (resp *types.ImageUpdateRes, err error) {

	if err := l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Select("file_name", "type", "status").
		Where("id = ?", req.Id).
		Updates(&models.Upload{
			FileName: req.FileName,
			Type:     req.Type,
			Status:   req.Status,
		}).
		Error; err != nil {
		return nil, err
	}

	if len(req.Tags) > 0 {
		if err := l.Remove(&models.UploadTag{}, req.Id); err != nil {
			return nil, err
		}

		var uploadTags []models.UploadTag

		for _, v := range req.Tags {
			uploadTags = append(
				uploadTags,
				models.UploadTag{
					UploadId: req.Id,
					TagId:    v,
				},
			)
		}

		if err := l.svcCtx.DB.
			Model(&models.UploadTag{}).
			Create(&uploadTags).
			Error; err != nil {
			return nil, err
		}
	}

	if len(req.Category) > 0 {
		if err := l.Remove(&models.UploadCategory{}, req.Id); err != nil {
			return nil, err
		}

		var uploadCategory []models.UploadCategory

		for _, v := range req.Category {
			uploadCategory = append(
				uploadCategory,
				models.UploadCategory{
					UploadId:   req.Id,
					CategoryId: v,
				},
			)
		}

		if err := l.svcCtx.DB.
			Model(&models.UploadCategory{}).
			Create(&uploadCategory).
			Error; err != nil {
			return nil, err
		}
	}

	if len(req.Recommend) > 0 {
		if err := l.Remove(&models.UploadRecommend{}, req.Id); err != nil {
			return nil, err
		}

		var uploadRecommend []models.UploadRecommend

		for _, v := range req.Recommend {
			uploadRecommend = append(
				uploadRecommend,
				models.UploadRecommend{
					UploadId:    req.Id,
					RecommendId: v,
				},
			)
		}

		if err := l.svcCtx.
			DB.
			Model(&models.UploadRecommend{}).
			Create(&uploadRecommend).
			Error; err != nil {
			return nil, err
		}
	}

	if len(req.Group) > 0 {
		if err := l.Remove(&models.UploadGroup{}, req.Id); err != nil {
			return nil, err
		}

		var uploadGroup []models.UploadGroup
		for _, v := range req.Group {
			uploadGroup = append(
				uploadGroup,
				models.UploadGroup{
					UploadId: req.Id,
					GroupId:  v,
				},
			)
		}

		if err := l.svcCtx.
			DB.
			Model(&models.UploadGroup{}).
			Create(&uploadGroup).
			Error; err != nil {
			return nil, err
		}
	} else {
		if err := l.Remove(&models.UploadGroup{}, req.Id); err != nil {
			return nil, err
		}
	}

	return &types.ImageUpdateRes{
		Base: types.Base{
			Code: 1,
			Msg:  "更新成功",
		},
	}, nil
}

func (l *ImageUpdateLogic) Remove(model any, id string) (err error) {
	if err = l.svcCtx.
		DB.
		Model(model).
		Where("upload_id = ?", id).
		Delete(model).
		Error; err != nil {
		return err
	}

	return nil
}
