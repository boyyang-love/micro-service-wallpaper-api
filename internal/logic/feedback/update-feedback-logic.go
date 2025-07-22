package feedback

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFeedbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFeedbackLogic {
	return &UpdateFeedbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateFeedbackLogic) UpdateFeedback(req *types.FeedbackUpdateReq) (resp *types.FeedbackUpdateRes, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.Feedback{}).
		Where("id = ?", req.Id).
		Update("status", req.Status).
		Error; err != nil {
		return nil, err
	}

	return &types.FeedbackUpdateRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
	}, nil
}
