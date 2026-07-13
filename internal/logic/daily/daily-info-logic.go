package daily

import (
	"context"
	"time"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type DailyInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDailyInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DailyInfoLogic {
	return &DailyInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DailyInfoLogic) DailyInfo(req *types.DailyInfoReq) (resp *types.DailyInfoRes, err error) {
	date := req.Date
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	var daily models.DailyWallpaper
	if err := l.svcCtx.DB.
		Where("date = ? AND status = 1", date).
		First(&daily).Error; err != nil {
		return &types.DailyInfoRes{
			Base: types.Base{Code: 1, Msg: "ok"},
			Data: nil,
		}, nil
	}

	var upload models.Upload
	if err := l.svcCtx.DB.
		Model(&models.Upload{}).
		Where("id = ?", daily.UploadId).
		First(&upload).Error; err != nil {
		return nil, err
	}

	var prevDate string
	l.svcCtx.DB.
		Model(&models.DailyWallpaper{}).
		Where("date < ? AND status = 1", date).
		Order("date DESC").
		Limit(1).
		Pluck("date", &prevDate)

	var nextDate string
	l.svcCtx.DB.
		Model(&models.DailyWallpaper{}).
		Where("date > ? AND status = 1", date).
		Order("date ASC").
		Limit(1).
		Pluck("date", &nextDate)

	var downloadsToday int64
	l.svcCtx.DB.
		Model(&models.Download{}).
		Where("download_id = ? AND DATE(FROM_UNIXTIME(created/1000)) = ?", daily.UploadId, date).
		Count(&downloadsToday)

	var tags []string
	type TagResult struct {
		Name string `json:"name"`
	}
	var tagResults []TagResult
	l.svcCtx.DB.
		Table("upload_tag ut").
		Select("t.name").
		Joins("LEFT JOIN tag t ON t.id = ut.tag_id").
		Where("ut.upload_id = ?", daily.UploadId).
		Find(&tagResults)
	for _, t := range tagResults {
		tags = append(tags, t.Name)
	}
	if tags == nil {
		tags = []string{}
	}

	return &types.DailyInfoRes{
		Base: types.Base{Code: 1, Msg: "ok"},
		Data: &types.DailyInfoData{
			Id:             upload.Id,
			FilePath:       upload.FilePath,
			FileName:       upload.FileName,
			W:              upload.W,
			H:              upload.H,
			Date:           daily.Date,
			Edition:        daily.Edition,
			Description:    daily.Description,
			DownloadsToday: int(downloadsToday),
			Tags:           tags,
			PrevDate:       prevDate,
			NextDate:       nextDate,
		},
	}, nil
}
