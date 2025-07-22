package download

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type Download struct {
	models.Download
	Upload Upload `gorm:"foreignkey:download_id;references:id"`
}

type Upload struct {
	models.Upload
}

func (d *Download) TableName() string {
	return "download"
}

func (l *Upload) TableName() string {
	return "upload"
}

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
	var count int64
	var download []Download
	records := make([]types.DownLoadUserListRecord, 0)

	DB := l.svcCtx.
		DB.
		Order("updated desc").
		Preload("Upload", func(db *gorm.DB) *gorm.DB {
			db = db.Select("id", "file_path", "file_name", "w", "h", "type")
			return db
		}).
		Model(&Download{})

	if req.Type != "" {
		DB = DB.Where("type = ?", req.Type)
	}

	if err =
		DB.
			Where("user_id = ?", userId).
			Offset((req.Page - 1) * req.Limit).
			Limit(req.Limit).
			Find(&download).
			Offset(-1).
			Count(&count).
			Error; err != nil {
		return nil, err
	}

	for _, v := range download {
		records = append(records, types.DownLoadUserListRecord{
			BaseTime: types.BaseTime{
				Created: v.Created,
				Updated: v.Updated,
			},
			Id:       v.Id,
			FileId:   v.Upload.Id,
			FilePath: v.Upload.FilePath,
			FileName: v.Upload.FileName,
			W:        v.Upload.W,
			H:        v.Upload.H,
			Type:     v.Upload.Type,
		})
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
			Records: records,
		},
	}, nil
}
