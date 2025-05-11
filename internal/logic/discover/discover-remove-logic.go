package discover

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/uploadclient"
	"strings"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DiscoverRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverRemoveLogic {
	return &DiscoverRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DiscoverRemoveLogic) DiscoverRemove(req *types.DiscoverRemoveReq) (resp *types.DiscoverRemoveRes, err error) {

	imageIds, err := l.GetImageIds(req.Id)
	if err != nil {
		return nil, err
	}

	err = l.DelUploadAndImages(imageIds)
	if err != nil {
		return nil, err
	}

	return &types.DiscoverRemoveRes{
		Base: types.Base{
			Code: 1,
			Msg:  "删除成功",
		},
	}, nil
}

func (l *DiscoverRemoveLogic) GetImageIds(id string) (imageIds []string, err error) {
	var imageIdsStr string
	if err = l.svcCtx.DB.
		Model(&models.Discover{}).
		Where("id = ?", id).
		Select("image_ids").
		Find(&imageIdsStr).
		Error; err != nil {
		return nil, err
	}

	if err := l.svcCtx.
		DB.
		Model(&models.Discover{}).
		Where("id = ?", id).
		Delete(&models.Discover{}).
		Error; err != nil {
		return nil, err
	}

	return strings.Split(imageIdsStr, ","), nil
}

func (l *DiscoverRemoveLogic) DelUploadAndImages(ids []string) error {
	var imagesPath []string
	var upload []models.Upload
	if err := l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Where("id in (?) and type = ?", ids, "DISCOVER").
		Select("file_path", "origin_file_path").
		Find(&upload).
		Error; err != nil {
		return err
	}

	for _, u := range upload {
		imagesPath = append(imagesPath, u.FilePath, u.OriginFilePath)
	}

	if err := l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Where("id in (?) and type = ?", ids, "DISCOVER").
		Delete(&models.Upload{}).
		Error; err != nil {
		return err
	}
	fmt.Println(imagesPath, ids)
	if _, err := l.svcCtx.UploadService.ImageDelete(l.ctx, &uploadclient.ImageDeleteReq{
		BucketName: "wallpaper",
		Paths:      imagesPath,
	}); err != nil {
		return err
	}

	return nil
}
