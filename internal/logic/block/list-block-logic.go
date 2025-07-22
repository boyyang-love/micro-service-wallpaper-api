package block

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListBlockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListBlockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBlockLogic {
	return &ListBlockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListBlockLogic) ListBlock(req *types.BlockListReq) (resp *types.BlockListRes, err error) {
	records := make([]types.BlockListInfo, 0)
	blocks := make([]models.Block, 0)
	var count int64
	if err = l.svcCtx.
		DB.
		Model(&models.Block{}).
		Where("type = ?", req.Type).
		Select("id", "created", "updated", "target_id", "type", "user_id").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&blocks).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	for _, block := range blocks {
		userInfo, dataInfo, err := l.ListBlockInfo(block.TargetId, req.Type)
		if err != nil {
			return nil, err
		}

		records = append(records, types.BlockListInfo{
			BaseTime: types.BaseTime{
				Created: block.Created,
				Updated: block.Updated,
			},
			Id:            block.Id,
			BlockUserInfo: userInfo,
			BlockDataInfo: dataInfo,
		})
	}

	return &types.BlockListRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.BlockListData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: records,
		},
	}, nil
}

func (l *ListBlockLogic) ListBlockInfo(targetId string, blockType string) (userInfo types.BlockUserInfo, dataInfo types.BlockDataInfo, err error) {
	if blockType == "user" {
		if err = l.svcCtx.
			DB.
			Model(&models.User{}).
			Where("id = ?", targetId).
			Select("id", "username", "avatar as user_avatar").
			First(&userInfo).
			Error; err != nil {
			return userInfo, dataInfo, err
		}

		return userInfo, dataInfo, nil
	}

	if blockType == "discover" {
		type Discover struct {
			Id       string `json:"id"`
			Title    string `json:"title"`
			ImageIds string `json:"image_ids"`
		}
		var discover Discover
		if err = l.svcCtx.
			DB.
			Model(&models.Discover{}).
			Where("id = ?", targetId).
			Select("id", "title", "image_ids").
			First(&discover).
			Error; err != nil {
			return userInfo, dataInfo, err
		}
		type Upload struct {
			Id             string `json:"id"`
			Type           string `json:"type"`
			OriginFilePath string `json:"origin_file_path"`
			FilePath       string `json:"file_path"`
			W              int    `json:"w"`
			H              int    `json:"h"`
		}
		var upload Upload
		if err = l.svcCtx.
			DB.
			Model(&models.Upload{}).
			Where("id = ?", discover.ImageIds).
			Select("id", "type", "origin_file_path", "file_path", "w", "h").
			First(&upload).
			Error; err != nil {
			return userInfo, dataInfo, err
		}

		return userInfo, types.BlockDataInfo{
			Id:             discover.Id,
			TargetId:       upload.Id,
			Title:          discover.Title,
			Type:           upload.Type,
			OriginFilePath: upload.OriginFilePath,
			FilePath:       upload.FilePath,
			W:              upload.W,
			H:              upload.H,
		}, err
	}

	return
}
