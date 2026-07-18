package post

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostReviewUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPostReviewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostReviewUpdateLogic {
	return &PostReviewUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostReviewUpdateLogic) PostReviewUpdate(req *types.PostReviewUpdateReq) (resp *types.PostReviewUpdateRes, err error) {
	if req.Id == "" {
		return &types.PostReviewUpdateRes{
			Base: types.Base{Code: 0, Msg: "id is required"},
		}, nil
	}

	if req.Status != 1 && req.Status != -1 {
		return &types.PostReviewUpdateRes{
			Base: types.Base{Code: 0, Msg: "status must be 1 or -1"},
		}, nil
	}

	updates := map[string]interface{}{
		"status": req.Status,
	}
	if req.Status == -1 && req.RejectReason != "" {
		updates["reject_reason"] = req.RejectReason
	}

	result := l.svcCtx.DB.Model(&models.Post{}).
		Where("id = ?", req.Id).
		Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return &types.PostReviewUpdateRes{
			Base: types.Base{Code: 0, Msg: "帖子不存在"},
		}, nil
	}

	return &types.PostReviewUpdateRes{
		Base: types.Base{Code: 1, Msg: "ok"},
	}, nil
}
