package comment

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentReviewListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentReviewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentReviewListLogic {
	return &CommentReviewListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentReviewListLogic) CommentReviewList(req *types.CommentReviewListReq) (resp *types.CommentReviewListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 10
	}

	var comments []models.Comment
	var count int64

	db := l.svcCtx.DB.Model(&models.Comment{})
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}

	if err = db.
		Order("created DESC").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&comments).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	records := make([]types.CommentReviewListItem, 0, len(comments))
	for _, c := range comments {
		// 用户信息
		var user models.User
		userInfo := types.CommentUserInfo{}
		if err := l.svcCtx.DB.Model(&models.User{}).
			Where("id = ?", c.UserId).
			Select("id", "username", "avatar").
			First(&user).Error; err == nil {
			userInfo = types.CommentUserInfo{
				Id:       user.Id,
				Username: user.Username,
				Avatar:   user.Avatar,
			}
		}

		// 帖子标题
		var post models.Post
		postTitle := ""
		if err := l.svcCtx.DB.Model(&models.Post{}).
			Where("id = ?", c.PostId).
			Select("title").
			First(&post).Error; err == nil {
			postTitle = post.Title
		}

		records = append(records, types.CommentReviewListItem{
			BaseTime: types.BaseTime{
				Created: c.Created,
				Updated: c.Updated,
			},
			Id:        c.Id,
			PostId:    c.PostId,
			PostTitle: postTitle,
			Content:   c.Content,
			Status:    c.Status,
			RejectReason: c.RejectReason,
			UserInfo:  userInfo,
		})
	}

	return &types.CommentReviewListRes{
		Base: types.Base{Code: 1, Msg: "ok"},
		Data: types.CommentReviewListData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: records,
		},
	}, nil
}
