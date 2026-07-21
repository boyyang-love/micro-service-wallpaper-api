// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package community

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProfileLogic {
	return &UserProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserProfileLogic) UserProfile(req *types.UserProfileReq) (resp *types.UserProfileRes, err error) {
	if req.UserId == "" {
		return &types.UserProfileRes{
			Base: types.Base{Code: 0, Msg: "user_id is required"},
		}, nil
	}

	// 查询用户信息
	var user models.User
	if err = l.svcCtx.DB.Model(&models.User{}).
		Where("id = ?", req.UserId).
		Select("id", "username", "avatar", "cover", "motto").
		First(&user).Error; err != nil {
		return &types.UserProfileRes{
			Base: types.Base{Code: 0, Msg: "user not found"},
		}, nil
	}

	// 统计上传数
	var postCount int64
	l.svcCtx.DB.Model(&models.Post{}).
		Where("user_id = ? AND status = ?", req.UserId, 1).
		Count(&postCount)

	// 统计获赞数（该用户上传的所有壁纸的点赞总数）
	var likeCount int64
	l.svcCtx.DB.Model(&models.Like{}).
		Where("upload_id IN (?) AND status = ?",
			l.svcCtx.DB.Model(&models.Post{}).Select("id").Where("user_id = ?", req.UserId),
			true,
		).Count(&likeCount)

	// 统计下载数（该用户上传的所有壁纸的 download 字段总和）
	var downloadCount int64
	l.svcCtx.DB.Model(&models.Upload{}).
		Where("user_id = ?", req.UserId).
		Select("COALESCE(SUM(download), 0)").
		Scan(&downloadCount)

	return &types.UserProfileRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.UserProfileData{
			UserInfo: types.UserProfileInfo{
				Id:       user.Id,
				Username: user.Username,
				Avatar:   user.Avatar,
				Cover:    user.Cover,
				Motto:    user.Motto,
			},
			Stats: types.UserProfileStats{
				UploadCount:   postCount,
				LikeCount:     likeCount,
				DownloadCount: downloadCount,
			},
		},
	}, nil
}
