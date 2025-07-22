package home

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HomeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHomeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomeListLogic {
	return &HomeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomeListLogic) HomeList(req *types.HomeListReq) (resp *types.HomeListRes, err error) {
	pcList, err := l.GetPc(req)
	if err != nil {
		return nil, err
	}

	moaList, err := l.GetMoa(req)
	if err != nil {
		return nil, err
	}

	hotPcList, err := l.GetHotPc(req)
	if err != nil {
		return nil, err
	}

	downloadMoa, err := l.GetDownloadMoa(req)
	if err != nil {
		return nil, err
	}

	avatar, err := l.GetAvatar(req)
	if err != nil {
		return nil, err
	}

	return &types.HomeListRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.HomeListResData{
			Pc:          pcList,
			Moa:         moaList,
			HotPc:       hotPcList,
			DownloadMoa: downloadMoa,
			HotAvatar:   avatar,
		},
	}, nil
}

func (l *HomeListLogic) GetPc(req *types.HomeListReq) (resp []types.HomeListInfo, err error) {
	list := make([]types.HomeListInfo, 0)
	if err = l.svcCtx.
		DB.
		Order("created desc").
		Model(&models.Upload{}).
		Select("id", "file_name", "file_path", "w", "h", "download", "view").
		Where("type = ? and status = 1", "PC").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&list).
		Offset(-1).
		Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (l *HomeListLogic) GetMoa(req *types.HomeListReq) (resp []types.HomeListInfo, err error) {
	list := make([]types.HomeListInfo, 0)
	if err = l.svcCtx.
		DB.
		Order("created desc").
		Model(&models.Upload{}).
		Select("id", "file_name", "file_path", "w", "h", "download", "view").
		Where("type = ? and status = 1", "MOA").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&list).
		Offset(-1).
		Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (l *HomeListLogic) GetHotPc(req *types.HomeListReq) (resp []types.HomeListInfo, err error) {
	list := make([]types.HomeListInfo, 0)
	if err = l.svcCtx.
		DB.
		Order("rand()").
		Model(&models.Upload{}).
		Select("id", "file_name", "file_path", "w", "h", "download", "view").
		Where("type = ? and status = 1", "PC").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&list).
		Offset(-1).
		Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (l *HomeListLogic) GetDownloadMoa(req *types.HomeListReq) (resp []types.HomeListInfo, err error) {
	list := make([]types.HomeListInfo, 0)
	if err = l.svcCtx.
		DB.
		Order("rand()").
		Model(&models.Upload{}).
		Select("id", "file_name", "file_path", "w", "h", "download", "view").
		Where("type = ? and status = 1", "MOA").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&list).
		Offset(-1).
		Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (l *HomeListLogic) GetAvatar(req *types.HomeListReq) (resp []types.HomeListInfo, err error) {
	list := make([]types.HomeListInfo, 0)
	if err = l.svcCtx.
		DB.
		Order("rand()").
		Model(&models.Upload{}).
		Select("id", "file_name", "file_path", "w", "h", "download", "view").
		Where("type = ? and status = 1", "AVATAR").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&list).
		Offset(-1).
		Error; err != nil {
		return nil, err
	}

	return list, nil
}
