package user

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/uploadclient"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoReq) (resp *types.UpdateUserInfoRes, err error) {
	var userId = fmt.Sprintf("%s", l.ctx.Value("Id"))

	if err = l.svcCtx.
		DB.
		Model(&models.User{}).
		Where("id = ?", userId).
		Updates(&models.User{
			Avatar: req.Avatar,
			Cover:  req.Cover,
		}).
		Error; err != nil {
		return nil, err
	}

	if err = l.RemoveUserAvatar(req.Avatar, userId); err != nil {
		return nil, err
	}

	return &types.UpdateUserInfoRes{
		Base: types.Base{
			Code: 1,
			Msg:  "更新成功",
		},
	}, nil
}

func (l *UpdateUserInfoLogic) RemoveUserAvatar(path string, userId string) error {
	var upload []models.Upload
	if err := l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Where("user_id = ? and type = ? and  file_path != ?", userId, "USERAVATAR", path).
		Find(&upload).
		Error; err != nil {
		return err
	}

	var willRemovePath []string
	var ids []string
	for _, u := range upload {
		willRemovePath = append(willRemovePath, u.FilePath, u.OriginFilePath)
		ids = append(ids, u.Id)
	}

	fmt.Println(ids, willRemovePath)

	if err := l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Delete(&models.Upload{}, ids).
		Error; err != nil {
		return err
	}

	_, err := l.svcCtx.UploadService.ImageDelete(l.ctx, &uploadclient.ImageDeleteReq{
		BucketName: "wallpaper",
		Paths:      willRemovePath,
	})
	if err != nil {
		return err
	}

	return nil
}
