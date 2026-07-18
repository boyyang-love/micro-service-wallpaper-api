package post

import (
	"context"
	"strings"

	"github.com/boyyang-love/micro-service-wallpaper-api/helper"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostFeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPostFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostFeedLogic {
	return &PostFeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostFeedLogic) PostFeed(req *types.PostFeedReq, authHeader string) (resp *types.PostFeedRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 20
	}

	var posts []models.Post
	var count int64

	db := l.svcCtx.DB.Model(&models.Post{}).Where("status = ?", 1)

	switch req.Sort {
	case "hot":
		db = db.Order("created DESC")
	default:
		db = db.Order("created DESC")
	}

	if err = db.
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&posts).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	// 从 Authorization header 解析用户 ID（可选）
	userId := l.parseUserId(authHeader)

	// 批量查询当前用户的点赞状态，避免 N+1 查询
	likedMap := make(map[string]bool)
	if userId != "" {
		var likedPosts []models.Like
		postIds := make([]string, 0, len(posts))
		for _, p := range posts {
			postIds = append(postIds, p.Id)
		}
		l.svcCtx.DB.Model(&models.Like{}).
			Where("upload_id IN (?) AND user_id = ? AND status = ?", postIds, userId, true).
			Find(&likedPosts)
		for _, like := range likedPosts {
			likedMap[like.UploadId] = true
		}
	}

	records := make([]types.PostFeedInfo, 0, len(posts))
	for _, p := range posts {
		images := l.getImages(p.ImageIds)
		userInfo := l.getUserInfo(p.UserId)
		likeCount := l.getLikeCount(p.Id)
		commentCount := l.getCommentCount(p.Id)

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
			IsLiked:      likedMap[p.Id],
		})
	}

	return &types.PostFeedRes{
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

func (l *PostFeedLogic) parseUserId(authHeader string) string {
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := helper.ParseToken(token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return ""
	}
	return claims.Id
}

func (l *PostFeedLogic) getImages(imageIds string) []types.PostImageInfo {
	if imageIds == "" {
		return []types.PostImageInfo{}
	}
	ids := strings.Split(imageIds, ",")
	var uploads []models.Upload
	l.svcCtx.DB.Model(&models.Upload{}).
		Select("id", "file_name", "file_path", "w", "h").
		Where("id IN (?)", ids).
		Find(&uploads)

	images := make([]types.PostImageInfo, 0, len(uploads))
	for _, u := range uploads {
		images = append(images, types.PostImageInfo{
			Id:       u.Id,
			FilePath: u.FilePath,
			FileName: u.FileName,
			W:        u.W,
			H:        u.H,
		})
	}
	return images
}

func (l *PostFeedLogic) getUserInfo(userId string) types.PostUserInfo {
	var user models.User
	if err := l.svcCtx.DB.Model(&models.User{}).
		Where("id = ?", userId).
		Select("id", "username", "avatar").
		First(&user).Error; err != nil {
		return types.PostUserInfo{}
	}
	return types.PostUserInfo{
		Id:       user.Id,
		Username: user.Username,
		Avatar:   user.Avatar,
	}
}

func (l *PostFeedLogic) getLikeCount(postId string) int64 {
	var count int64
	l.svcCtx.DB.Model(&models.Like{}).
		Where("upload_id = ? AND status = ?", postId, true).
		Count(&count)
	return count
}

func (l *PostFeedLogic) getCommentCount(postId string) int64 {
	var count int64
	l.svcCtx.DB.Model(&models.Comment{}).
		Where("post_id = ? AND status = ?", postId, 1).
		Count(&count)
	return count
}

func (l *PostFeedLogic) checkIsLiked(postId, userId string) bool {
	var like models.Like
	err := l.svcCtx.DB.Model(&models.Like{}).
		Where("upload_id = ? AND user_id = ? AND status = ?", postId, userId, true).
		First(&like).Error
	return err == nil
}
