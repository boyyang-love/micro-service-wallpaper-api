package album

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailAlbumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailAlbumLogic {
	return &DetailAlbumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailAlbumLogic) DetailAlbum(req *types.AlbumDetailReq) (resp *types.AlbumDetailRes, err error) {
	count, ids, err := l.ids(req)
	if err != nil {
		return nil, err
	}
	var upload []models.Upload
	records := make([]types.AlbumImageInfo, 0)

	if err = l.svcCtx.
		DB.
		Order("created desc").
		Model(&models.Upload{}).
		Where("id in (?)", ids).
		Select("id", "file_name", "file_path").
		Find(&upload).
		Error; err != nil {
		return nil, err
	}
	for _, u := range upload {
		records = append(records, types.AlbumImageInfo{
			Id:       u.Id,
			FileName: u.FileName,
			FilePath: u.FilePath,
		})
	}

	return &types.AlbumDetailRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.AlbumDetailData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: records,
		},
	}, nil
}

func (l *DetailAlbumLogic) ids(req *types.AlbumDetailReq) (count int64, ids []string, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.UploadAlbum{}).
		Where("album_id = ?", req.Id).
		Select("upload_id").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&ids).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return count, nil, err
	}

	return count, ids, nil
}
