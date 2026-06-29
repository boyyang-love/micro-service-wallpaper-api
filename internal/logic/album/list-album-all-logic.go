package album

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAlbumAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAlbumAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAlbumAllLogic {
	return &ListAlbumAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAlbumAllLogic) ListAlbumAll(req *types.AlbumListReq) (resp *types.AlbumListRes, err error) {
	records := make([]types.Album, 0)
	albums := make([]models.Album, 0)
	var count int64
	if err = l.svcCtx.
		DB.
		Order("created desc").
		Model(&models.Album{}).
		Offset((req.Page - 1) * req.Limit).
		Find(&albums).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	for _, album := range albums {
		imagePath, _ := l.getImagePathById(album.Cover)
		records = append(records, types.Album{
			Id:        album.Id,
			Name:      album.Name,
			Desc:      album.Desc,
			Cover:     album.Cover,
			CoverPath: imagePath,
			Status:    album.Status,
		})
	}

	return &types.AlbumListRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.AlbumListData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: records,
		},
	}, nil
}

func (l *ListAlbumAllLogic) getImagePathById(id string) (imagePath string, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Where("id = ?", id).
		Select("file_path").
		Scan(&imagePath).
		Error; err != nil {
		return "", err
	}

	return imagePath, nil
}
