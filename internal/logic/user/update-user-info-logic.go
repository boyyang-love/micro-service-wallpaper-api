package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/uploadclient"
	"gorm.io/gorm"

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

	// 先查询旧的用户信息
	var oldUser models.User
	if err = l.svcCtx.DB.Model(&models.User{}).Where("id = ?", userId).First(&oldUser).Error; err != nil {
		return nil, err
	}

	// 更新用户信息
	if err = l.svcCtx.DB.Model(&models.User{}).
		Where("id = ?", userId).
		Updates(&models.User{
			Avatar:   req.Avatar,
			Cover:    req.Cover,
			Username: req.Username,
			Motto:    req.Motto,
		}).Error; err != nil {
		return nil, err
	}

	// 如果更新了头像，删除旧头像
	if req.Avatar != "" && oldUser.Avatar != "" && req.Avatar != oldUser.Avatar {
		l.RemoveOldImage(userId, "AVATAR", oldUser.Avatar)
	}

	// 如果更新了封面，删除旧封面
	if req.Cover != "" && oldUser.Cover != "" && req.Cover != oldUser.Cover {
		l.RemoveOldImage(userId, "COVER", oldUser.Cover)
	}

	return &types.UpdateUserInfoRes{
		Base: types.Base{
			Code: 1,
			Msg:  "更新成功",
		},
	}, nil
}

// RemoveOldImage 删除用户旧的头像/封面图片
func (l *UpdateUserInfoLogic) RemoveOldImage(userId string, imageType string, oldPath string) {
	// 跳过外部链接（QQ头像等）
	if oldPath == "" || len(oldPath) > 4 && oldPath[:4] == "http" {
		return
	}

	var upload models.Upload
	if err := l.svcCtx.DB.Model(&models.Upload{}).
		Where("user_id = ? AND type = ? AND file_path = ?", userId, imageType, oldPath).
		First(&upload).Error; err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			return
		}
		l.Logger.Error("查询旧图片失败:", err)
		return
	}

	// 删除数据库记录
	if err := l.svcCtx.DB.Delete(&upload).Error; err != nil {
		l.Logger.Error("删除旧图片记录失败:", err)
		return
	}

	// 删除存储文件
	var paths []string
	if upload.FilePath != "" {
		paths = append(paths, upload.FilePath)
	}
	if upload.OriginFilePath != "" {
		paths = append(paths, upload.OriginFilePath)
	}

	if len(paths) > 0 {
		_, err := l.svcCtx.UploadService.CosDelete(l.ctx, &uploadclient.ImageDeleteReq{
			BucketName: "wallpaper",
			Paths:      paths,
		})
		if err != nil {
			l.Logger.Error("删除旧图片文件失败:", err)
		}
	}
}
