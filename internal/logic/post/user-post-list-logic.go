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

type UserPostListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserPostListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserPostListLogic {
	return &UserPostListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserPostListLogic) UserPostList(req *types.UserPostListReq) (resp *types.UserPostListRes, err error) {
	if req.UserId == "" {
		return &types.UserPostListRes{
			Base: types.Base{Code: 0, Msg: "user_id is required"},
		}, nil
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 20
	}

	var posts []models.Post
	var count int64

	if err = l.svcCtx.DB.Model(&models.Post{}).
		Where("user_id = ?", req.UserId).
		Order("created DESC").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&posts).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	userId := fmt.Sprintf("%s", l.ctx.Value("Id"))

	records := make([]types.PostFeedInfo, 0, len(posts))
	for _, p := range posts {
		var images []types.PostImageInfo
		if p.ImageIds != "" {
			ids := strings.Split(p.ImageIds, ",")
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

		var user models.User
		l.svcCtx.DB.Model(&models.User{}).
			Where("id = ?", p.UserId).
			Select("id", "username", "avatar").
			First(&user)
		userInfo := types.PostUserInfo{
			Id:       user.Id,
			Username: user.Username,
			Avatar:   user.Avatar,
		}

		var likeCount int64
		l.svcCtx.DB.Model(&models.Like{}).
			Where("upload_id = ? AND status = ?", p.Id, true).
			Count(&likeCount)

		var commentCount int64
		l.svcCtx.DB.Model(&models.Comment{}).
			Where("post_id = ? AND status = ?", p.Id, 1).
			Count(&commentCount)

		isLiked := false
		if userId != "" && userId != "<nil>" {
			var like models.Like
			if l.svcCtx.DB.Model(&models.Like{}).
				Where("upload_id = ? AND user_id = ? AND status = ?", p.Id, userId, true).
				First(&like).Error == nil {
				isLiked = true
			}
		}

		records = append(records, types.PostFeedInfo{
			BaseTime: types.BaseTime{
				Created: p.Created,
				Updated: p.Updated,
			},
			Id:           p.Id,
			Title:        p.Title,
			Content:      p.Content,
			Images:       images,
			UserInfo:     userInfo,
			LikeCount:    likeCount,
			CommentCount: commentCount,
			Status:        p.Status,
			RejectReason:  p.RejectReason,
			IsLiked:      isLiked,
		})
	}

	return &types.UserPostListRes{
		Base: types.Base{Code: 1, Msg: "ok"},
		Data: types.PostFeedData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: records,
		},
	}, nil
}
