package like

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveAllLikeRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveAllLikeRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveAllLikeRecordLogic {
	return &RemoveAllLikeRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveAllLikeRecordLogic) RemoveAllLikeRecord() (resp *types.RemoveAllLikeRecordRes, err error) {
	userid := fmt.Sprintf("%s", l.ctx.Value("Id"))

	if err = l.svcCtx.
		DB.
		Model(&models.Like{}).
		Where("userid = ?", userid).
		Delete(&models.Like{}).
		Error; err != nil {
		return nil, err
	}

	return &types.RemoveAllLikeRecordRes{
		Base: types.Base{
			Code: 1,
			Msg:  "删除收藏成功",
		},
	}, nil
}
