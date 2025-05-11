package download

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddDownloadRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddDownloadRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddDownloadRecordLogic {
	return &AddDownloadRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddDownloadRecordLogic) AddDownloadRecord(req *types.AddDownloadRecordReq) (resp *types.AddDownloadRecordRes, err error) {
	userid := fmt.Sprintf("%s", l.ctx.Value("Id"))

	if err := l.svcCtx.
		DB.
		Model(&models.Download{}).
		Create(&models.Download{
			DownloadId: req.DownloadId,
			UserId:     userid,
		}).Error; err != nil {
		return nil, err
	}

	return &types.AddDownloadRecordRes{
		Base: types.Base{
			Code: 1,
			Msg:  "新增记录成功",
		},
	}, nil
}
