package daily

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type DailyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDailyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DailyListLogic {
	return &DailyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DailyListLogic) DailyList(req *types.DailyListReq) (resp *types.DailyListRes, err error) {
	var dailies []models.DailyWallpaper
	var count int64

	query := l.svcCtx.DB.Model(&models.DailyWallpaper{}).Where("status = 1")

	if err := query.Count(&count).Error; err != nil {
		return nil, err
	}

	if err := query.
		Order("date DESC").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&dailies).Error; err != nil {
		return nil, err
	}

	var records []types.DailyListInfo
	for _, d := range dailies {
		var upload models.Upload
		if err := l.svcCtx.DB.
			Model(&models.Upload{}).
			Where("id = ?", d.UploadId).
			First(&upload).Error; err != nil {
			continue
		}

		records = append(records, types.DailyListInfo{
			BaseTime: types.BaseTime{
				Created: d.Created,
				Updated: d.Updated,
			},
			Id:       upload.Id,
			DailyId:  d.Id,
			FilePath: upload.FilePath,
			FileName: upload.FileName,
			Date:     d.Date,
			Edition:  d.Edition,
			W:        upload.W,
			H:        upload.H,
		})
	}

	if records == nil {
		records = []types.DailyListInfo{}
	}

	return &types.DailyListRes{
		Base: types.Base{Code: 1, Msg: "ok"},
		Data: types.DailyListData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: records,
		},
	}, nil
}
