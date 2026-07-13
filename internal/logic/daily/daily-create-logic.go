package daily

import (
	"context"
	"time"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type DailyCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDailyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DailyCreateLogic {
	return &DailyCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DailyCreateLogic) DailyCreate(req *types.DailyCreateReq) (resp *types.DailyCreateRes, err error) {
	date := req.Date
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	var count int64
	l.svcCtx.DB.Model(&models.DailyWallpaper{}).
		Where("date = ?", date).
		Count(&count)
	if count > 0 {
		return &types.DailyCreateRes{
			Base: types.Base{Code: 1, Msg: "该日期已有每日壁纸"},
		}, nil
	}

	var maxEdition int
	l.svcCtx.DB.Model(&models.DailyWallpaper{}).
		Select("COALESCE(MAX(edition), 0)").
		Scan(&maxEdition)

	daily := models.DailyWallpaper{
		UploadId:    req.UploadId,
		Date:        date,
		Description: req.Description,
		Edition:     maxEdition + 1,
		Status:      1,
	}

	if err := l.svcCtx.DB.Create(&daily).Error; err != nil {
		return nil, err
	}

	return &types.DailyCreateRes{
		Base: types.Base{Code: 1, Msg: "ok"},
	}, nil
}
