package comment

import (
	"context"
	"fmt"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentCreateLogic {
	return &CommentCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentCreateLogic) CommentCreate(req *types.CommentCreateReq) (resp *types.CommentCreateRes, err error) {
	userId := fmt.Sprintf("%s", l.ctx.Value("Id"))
	if userId == "" || userId == "<nil>" {
		return &types.CommentCreateRes{
			Base: types.Base{Code: 0, Msg: "请先登录"},
		}, nil
	}

	if req.PostId == "" {
		return &types.CommentCreateRes{
			Base: types.Base{Code: 0, Msg: "post_id is required"},
		}, nil
	}

	if req.Content == "" {
		return &types.CommentCreateRes{
			Base: types.Base{Code: 0, Msg: "评论内容不能为空"},
		}, nil
	}

	comment := models.Comment{
		PostId:  req.PostId,
		UserId:  userId,
		Content: req.Content,
		Status:  2,
	}

	if err = l.svcCtx.DB.Create(&comment).Error; err != nil {
		return nil, err
	}

	return &types.CommentCreateRes{
		Base: types.Base{Code: 1, Msg: "ok"},
		Data: types.CommentCreateData{Id: comment.Id},
	}, nil
}
