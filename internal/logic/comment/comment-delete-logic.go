package comment

import (
	"context"
	"fmt"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentDeleteLogic {
	return &CommentDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentDeleteLogic) CommentDelete(req *types.CommentDeleteReq) (resp *types.CommentDeleteRes, err error) {
	userId := fmt.Sprintf("%s", l.ctx.Value("Id"))
	if userId == "" || userId == "<nil>" {
		return &types.CommentDeleteRes{
			Base: types.Base{Code: 0, Msg: "请先登录"},
		}, nil
	}

	// 只能删除自己的评论，硬删除
	result := l.svcCtx.DB.
		Where("id = ? AND user_id = ?", req.Id, userId).
		Delete(&models.Comment{})

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return &types.CommentDeleteRes{
			Base: types.Base{Code: 0, Msg: "评论不存在或无权删除"},
		}, nil
	}

	return &types.CommentDeleteRes{
		Base: types.Base{Code: 1, Msg: "ok"},
	}, nil
}
