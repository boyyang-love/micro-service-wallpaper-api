package comment

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req *types.CommentListReq) (resp *types.CommentListRes, err error) {
	if req.PostId == "" {
		return &types.CommentListRes{
			Base: types.Base{Code: 0, Msg: "post_id is required"},
		}, nil
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 50
	}

	var comments []models.Comment
	var count int64

	if err = l.svcCtx.DB.Model(&models.Comment{}).
		Where("post_id = ? AND status = ?", req.PostId, 1).
		Order("created DESC").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&comments).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	records := make([]types.CommentInfo, 0, len(comments))
	for _, c := range comments {
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

		records = append(records, types.CommentInfo{
			BaseTime: types.BaseTime{
				Created: c.Created,
				Updated: c.Updated,
			},
			Id:       c.Id,
			Content:  c.Content,
			UserInfo: userInfo,
		})
	}

	return &types.CommentListRes{
		Base: types.Base{Code: 1, Msg: "ok"},
		Data: types.CommentListData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: records,
		},
	}, nil
}
