package feedback

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFeedbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFeedbackLogic {
	return &CreateFeedbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateFeedbackLogic) CreateFeedback(req *types.FeedbackCreateReq) (resp *types.FeedbackCreateRes, err error) {
	userid := fmt.Sprintf("%s", l.ctx.Value("Id"))
	if err = l.svcCtx.DB.Model(&models.Feedback{}).Create(&models.Feedback{
		FeedbackId: req.FeedbackId,
		Content:    req.Content,
		Type:       req.Type,
		Status:     0,
		UserId:     userid,
	}).Error; err != nil {
		return nil, err
	}

	return &types.FeedbackCreateRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
	}, nil
}
