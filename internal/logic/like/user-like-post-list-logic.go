package like

import (
	"context"
	"fmt"
	"strings"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLikePostListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLikePostListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLikePostListLogic {
	return &UserLikePostListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLikePostListLogic) UserLikePostList(req *types.UserLikePostListReq) (resp *types.UserLikePostListRes, err error) {
	userId := fmt.Sprintf("%s", l.ctx.Value("Id"))
	if userId == "" || userId == "<nil>" {
		return &types.UserLikePostListRes{
			Base: types.Base{Code: 0, Msg: "请先登录"},
		}, nil
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 20
	}

	// 查询用户收藏的帖子（type=POST, status=true）
	var likes []models.Like
	var count int64

	if err = l.svcCtx.DB.Model(&models.Like{}).
		Where("user_id = ? AND type = ? AND status = ?", userId, "POST", true).
		Order("updated DESC").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&likes).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	records := make([]types.UserLikePostListItem, 0, len(likes))
	for _, like := range likes {
		// 查询帖子信息
		var post models.Post
		if err := l.svcCtx.DB.Model(&models.Post{}).
			Where("id = ? AND status = ?", like.UploadId, 1).
			First(&post).Error; err != nil {
			continue // 帖子已删除则跳过
		}

		// 获取封面图（第一张图片路径）
		coverPath := ""
		if post.ImageIds != "" {
			ids := strings.Split(post.ImageIds, ",")
			if len(ids) > 0 {
				var upload models.Upload
				if err := l.svcCtx.DB.Model(&models.Upload{}).
					Where("id = ?", ids[0]).
					Select("file_path").
					First(&upload).Error; err == nil {
					coverPath = upload.FilePath
				}
			}
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

		records = append(records, types.UserLikePostListItem{
			BaseTime: types.BaseTime{
				Created: post.Created,
				Updated: post.Updated,
			},
			Id:           like.Id,
			PostId:       post.Id,
			Title:        post.Title,
			Content:      post.Content,
			CoverPath:    coverPath,
			LikeCount:    likeCount,
			CommentCount: commentCount,
		})
	}

	return &types.UserLikePostListRes{
		Base: types.Base{Code: 1, Msg: "ok"},
		Data: types.UserLikePostListData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: records,
		},
	}, nil
}
