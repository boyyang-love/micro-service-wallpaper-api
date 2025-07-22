package download

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveAllDownloadRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveAllDownloadRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveAllDownloadRecordLogic {
	return &RemoveAllDownloadRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveAllDownloadRecordLogic) RemoveAllDownloadRecord() (resp *types.RemoveAllDownloadRecordRes, err error) {
	userid := fmt.Sprintf("%s", l.ctx.Value("Id"))

	if err = l.svcCtx.DB.
		Model(&models.Download{}).
		Where("userid = ?", userid).
		Delete(&models.Download{}).
		Error; err != nil {
		return nil, err
	}

	return &types.RemoveAllDownloadRecordRes{
		Base: types.Base{
			Code: 1,
			Msg:  "删除下载记录成功",
		},
	}, nil
}
