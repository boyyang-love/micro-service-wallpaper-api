package like

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeCreateLogic {
	return &LikeCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeCreateLogic) LikeCreate(req *types.LikeCreateOrUpdateReq) (resp *types.LikeCreateOrUpdateRes, err error) {
	userid := fmt.Sprintf("%s", l.ctx.Value("Id"))

	if l.svcCtx.
		DB.
		Model(&models.Like{}).
		Where("upload_id = ? and user_id = ?", req.UploadId, userid).
		Update("status", req.Status).
		RowsAffected == 0 {
		if err := l.svcCtx.DB.Create(&models.Like{
			UploadId: req.UploadId,
			UserId:   userid,
			Status:   req.Status,
		}).Error; err != nil {
			return nil, err
		}
	}
	return &types.LikeCreateOrUpdateRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
	}, nil
}
