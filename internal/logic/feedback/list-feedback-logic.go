package feedback

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFeedbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFeedbackLogic {
	return &ListFeedbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListFeedbackLogic) ListFeedback(req *types.FeedbackListReq) (resp *types.FeedbackListRes, err error) {

	records := make([]types.FeedbackInfo, 0)
	var count int64
	DB := l.svcCtx.DB.Model(&models.Feedback{})

	if req.Status != "" {
		DB = DB.Where("status = ?", req.Status)
	}

	if err = DB.
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&records).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	return &types.FeedbackListRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.FeedbackListData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: records,
		},
	}, nil
}
