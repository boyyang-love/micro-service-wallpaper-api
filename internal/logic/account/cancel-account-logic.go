package account

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

type CancelAccountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelAccountLogic {
	return &CancelAccountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelAccountLogic) CancelAccount(req *types.CancelAccountReq) (resp *types.CancelAccountRes, err error) {
	err = l.RemoveLike()
	if err != nil {
		return nil, err
	}

	err = l.RemoveDownload()
	if err != nil {
		return nil, err
	}

	err = l.RemoveDiscover()
	if err != nil {
		return nil, err
	}

	err = l.RemoveUser()
	if err != nil {
		return nil, err
	}

	err = l.RemoveBlock()
	if err != nil {
		return nil, err
	}

	return &types.CancelAccountRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
	}, nil
}

func (l *CancelAccountLogic) RemoveLike() error {
	userid := fmt.Sprintf("%s", l.ctx.Value("Id"))

	if err := l.svcCtx.
		DB.
		Model(&models.Like{}).
		Where("user_id = ?", userid).
		Delete(&models.Like{}).
		Error; err != nil {
		return err
	}

	return nil
}

func (l *CancelAccountLogic) RemoveDownload() error {
	userid := fmt.Sprintf("%s", l.ctx.Value("Id"))

	if err := l.svcCtx.
		DB.
		Model(&models.Download{}).
		Where("user_id = ?", userid).
		Delete(&models.Download{}).
		Error; err != nil {
		return err
	}

	return nil
}

func (l *CancelAccountLogic) RemoveDiscover() error {
	userid := fmt.Sprintf("%s", l.ctx.Value("Id"))
	imageIds, err := l.GetImageIds()
	if err != nil {
		return err
	}

	err = l.DelUploadAndImages(imageIds)
	if err != nil {
		return err
	}

	if err := l.svcCtx.
		DB.
		Model(&models.Discover{}).
		Where("user_id = ?", userid).
		Delete(&models.Discover{}).
		Error; err != nil {
		return err
	}

	return nil
}

func (l *CancelAccountLogic) RemoveUserAvatar() error {
	userid := fmt.Sprintf("%s", l.ctx.Value("Id"))

	if err := l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Where("user_id = ? and type = ?", userid, "USERAVATAR").
		Delete(&models.Upload{}).
		Error; err != nil {
		return err
	}

	return nil
}

func (l *CancelAccountLogic) RemoveUser() error {
	userid := fmt.Sprintf("%s", l.ctx.Value("Id"))

	if err := l.svcCtx.
		DB.
		Model(&models.User{}).
		Where("id = ?", userid).
		Delete(&models.User{}).
		Error; err != nil {
		return err
	}

	return nil
}

func (l *CancelAccountLogic) RemoveBlock() error {
	userid := fmt.Sprintf("%s", l.ctx.Value("Id"))

	if err := l.svcCtx.
		DB.
		Model(&models.Block{}).
		Where("user_id = ?", userid).
		Delete(&models.Block{}).
		Error; err != nil {
		return err
	}

	return nil
}

func (l *CancelAccountLogic) GetImageIds() (imageIds []string, err error) {
	userid := fmt.Sprintf("%s", l.ctx.Value("Id"))
	var discoverImageIds []string
	var userAvatarImageIds []string
	if err = l.svcCtx.DB.
		Model(&models.Discover{}).
		Where("user_id = ?", userid).
		Select("image_ids").
		Find(&discoverImageIds).
		Error; err != nil {
		return nil, err
	}

	if err = l.svcCtx.DB.
		Model(&models.Upload{}).
		Where("user_id = ? and type = ?", userid, "USERAVATAR").
		Select("id").
		Find(&userAvatarImageIds).
		Error; err != nil {
		return nil, err
	}

	return strings.Split(strings.Join(append(discoverImageIds, userAvatarImageIds...), ","), ","), nil
}

func (l *CancelAccountLogic) DelUploadAndImages(ids []string) error {
	var imagesPath []string
	var upload []models.Upload
	if err := l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Where("id in (?) and type in (?)", ids, []string{"DISCOVER", "USERAVATAR"}).
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
		Where("id in (?)", ids).
		Delete(&models.Upload{}).
		Error; err != nil {
		return err
	}

	if _, err := l.svcCtx.UploadService.CosDelete(l.ctx, &uploadclient.ImageDeleteReq{
		BucketName: "wallpaper",
		Paths:      imagesPath,
	}); err != nil {
		return err
	}

	return nil
}
