package post

import (
	"context"
	"fmt"
	"strings"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPostDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostDetailLogic {
	return &PostDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostDetailLogic) PostDetail(req *types.PostDetailReq) (resp *types.PostDetailRes, err error) {
	if req.Id == "" {
		return &types.PostDetailRes{
			Base: types.Base{Code: 0, Msg: "id is required"},
		}, nil
	}

	var post models.Post
	if err = l.svcCtx.DB.Model(&models.Post{}).
		Where("id = ? AND status = ?", req.Id, 1).
		First(&post).Error; err != nil {
		return &types.PostDetailRes{
			Base: types.Base{Code: 0, Msg: "帖子不存在"},
		}, nil
	}

	// 图片
	var images []types.PostImageInfo
	if post.ImageIds != "" {
		ids := strings.Split(post.ImageIds, ",")
		var uploads []models.Upload
		l.svcCtx.DB.Model(&models.Upload{}).
			Select("id", "file_name", "file_path", "w", "h").
			Where("id IN (?)", ids).
			Find(&uploads)
		for _, u := range uploads {
			images = append(images, types.PostImageInfo{
				Id:       u.Id,
				FilePath: u.FilePath,
				FileName: u.FileName,
				W:        u.W,
				H:        u.H,
			})
		}
	}

	// 用户信息
	var user models.User
	l.svcCtx.DB.Model(&models.User{}).
		Where("id = ?", post.UserId).
		Select("id", "username", "avatar").
		First(&user)
	userInfo := types.PostUserInfo{
		Id:       user.Id,
		Username: user.Username,
		Avatar:   user.Avatar,
	}

	// 点赞数
	var likeCount int64
	l.svcCtx.DB.Model(&models.Like{}).
		Where("upload_id = ? AND status = ?", post.Id, true).
		Count(&likeCount)

	// 评论数
	var commentCount int64
	l.svcCtx.DB.Model(&models.Comment{}).
		Where("post_id = ? AND status = ?", post.Id, 1).
		Count(&commentCount)

	// 是否已点赞
	isLiked := false
	userId := fmt.Sprintf("%s", l.ctx.Value("Id"))
	if userId != "" && userId != "<nil>" {
		var like models.Like
		if l.svcCtx.DB.Model(&models.Like{}).
			Where("upload_id = ? AND user_id = ? AND status = ?", post.Id, userId, true).
			First(&like).Error == nil {
			isLiked = true
		}
	}

	return &types.PostDetailRes{
		Base: types.Base{Code: 1, Msg: "ok"},
		Data: types.PostFeedInfo{
			BaseTime: types.BaseTime{
				Created: post.Created,
				Updated: post.Updated,
			},
			Id:           post.Id,
			Title:        post.Title,
			Content:      post.Content,
			Images:       images,
			UserInfo:     userInfo,
			LikeCount:    likeCount,
			CommentCount: commentCount,
			IsLiked:      isLiked,
		},
	}, nil
}
