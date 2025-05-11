package download

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDownloadListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDownloadListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDownloadListLogic {
	return &UserDownloadListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDownloadListLogic) UserDownloadList(req *types.DownlaodUserListReq) (resp *types.DownlaodUserListRes, err error) {

	var userId = fmt.Sprintf("%s", l.ctx.Value("Id"))
	var records []types.DownLoadUserListRecord
	var count int64

	uploadIds, err := l.DownloadIds(userId)
	if err != nil {
		return nil, err
	}

	DB := l.svcCtx.
		DB.
		Order("created desc").
		Model(&models.Upload{})

	if req.Type != "" {
		DB = DB.Where("type = ?", req.Type)
	}

	if err = DB.
		Where("id in (?)", uploadIds).
		Select("created", "updated", "id", "file_path", "file_name", "w", "h", "type").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&records).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	var sortedRecords []types.DownLoadUserListRecord

	for _, uploadId := range uploadIds {
		for _, record := range records {
			if uploadId == record.Id {
				sortedRecords = append(sortedRecords, record)
			}
		}
	}

	return &types.DownlaodUserListRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.DownloadUserListData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: sortedRecords,
		},
	}, nil
}

func (l *UserDownloadListLogic) DownloadIds(userId string) (ids []string, err error) {
	if err = l.svcCtx.
		DB.
		Order("created desc").
		Model(&models.Download{}).
		Distinct("download_id").
		Select("download_id").
		Where("user_id = ?", userId).
		Find(&ids).
		Error; err != nil {
		return ids, err
	}

	return ids, nil
}
