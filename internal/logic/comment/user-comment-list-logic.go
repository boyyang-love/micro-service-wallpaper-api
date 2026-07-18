package comment

import (
	"context"
	"fmt"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCommentListLogic {
	return &UserCommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCommentListLogic) UserCommentList(req *types.UserCommentListReq) (resp *types.UserCommentListRes, err error) {
	userId := fmt.Sprintf("%s", l.ctx.Value("Id"))
	if userId == "" || userId == "<nil>" {
		return &types.UserCommentListRes{
			Base: types.Base{Code: 0, Msg: "请先登录"},
		}, nil
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 20
	}

	var comments []models.Comment
	var count int64

	if err = l.svcCtx.DB.Model(&models.Comment{}).
		Where("user_id = ?", userId).
		Order("created DESC").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&comments).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	records := make([]types.UserCommentListItem, 0, len(comments))
	for _, c := range comments {
		// 查询帖子标题
		var post models.Post
		postTitle := ""
		if err := l.svcCtx.DB.Model(&models.Post{}).
			Where("id = ?", c.PostId).
			Select("title").
			First(&post).Error; err == nil {
			postTitle = post.Title
		}

		records = append(records, types.UserCommentListItem{
			BaseTime: types.BaseTime{
				Created: c.Created,
				Updated: c.Updated,
			},
			Id:           c.Id,
			PostId:       c.PostId,
			PostTitle:    postTitle,
			Content:      c.Content,
			Status:       c.Status,
			RejectReason: c.RejectReason,
		})
	}

	return &types.UserCommentListRes{
		Base: types.Base{Code: 1, Msg: "ok"},
		Data: types.UserCommentListData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: records,
		},
	}, nil
}
