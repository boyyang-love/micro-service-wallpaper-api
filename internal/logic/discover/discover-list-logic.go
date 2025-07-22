package discover

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/jinzhu/copier"
	"strings"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DiscoverListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverListLogic {
	return &DiscoverListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DiscoverListLogic) DiscoverList(req *types.DiscoverListReq) (resp *types.DiscoverListRes, err error) {

	var discover []models.Discover
	var records = make([]types.DiscoverListInfo, 0)
	var count int64

	DB := l.svcCtx.
		DB.
		Debug().
		Model(&models.Discover{})

	if req.Status != 0 {
		DB = DB.Where("status = ?", req.Status)
	}

	if req.UserId != "" {
		//DB = DB.Or("user_id = ?", req.UserId)

		userIds, discoverIds, err := l.GetBlockInfo(req.UserId)
		if err != nil {
			return nil, err
		}
		if len(discoverIds) > 0 {
			DB = DB.Where("id not in (?)", discoverIds)
		}

		if len(userIds) > 0 {
			DB = DB.Where("user_id not in (?)", userIds)
		}
	}

	if req.Sort != "" {
		DB = DB.Order(req.Sort)
	}

	if err := DB.
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&discover).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	for _, d := range discover {
		ids := strings.Split(d.ImageIds, ",")
		images, err := l.GetImages(ids)
		if err != nil {
			return nil, err
		}

		userInfo, err := l.GetUserInfo(d.UserId)
		if err != nil {
			return nil, err
		}

		records = append(records, types.DiscoverListInfo{
			BaseTime: types.BaseTime{
				Created: d.Created,
				Updated: d.Updated,
			},
			Id:       d.Id,
			Title:    d.Title,
			Subtitle: d.Subtitle,
			Status:   d.Status,
			UserInfo: *userInfo,
			Images:   images,
		})
	}

	return &types.DiscoverListRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.DiscoverListData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: records,
		},
	}, nil
}

func (l *DiscoverListLogic) GetImages(ids []string) (images []types.DiscoverListImages, err error) {
	var upload []models.Upload
	if err := l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Select("id", "file_name", "file_path", "w", "h", "type").
		Where("id in (?)", ids).
		Find(&upload).
		Error; err != nil {
		return nil, err
	}
	err = copier.Copy(&images, &upload)
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (l *DiscoverListLogic) GetUserInfo(userId string) (userInfo *types.DiscoverUserInfo, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.User{}).
		Where("id = ?", userId).
		Select("id", "username", "avatar").
		First(&userInfo).
		Error; err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (l *DiscoverListLogic) GetBlockInfo(userId string) (userIds []string, discoverIds []string, err error) {
	var block []models.Block
	if err = l.svcCtx.
		DB.
		Model(&models.Block{}).
		Select("type", "target_id", "user_id").
		Where("user_id = ?", userId).
		Find(&block).
		Error; err != nil {
		return nil, nil, err
	}

	for _, b := range block {
		if b.Type == "user" {
			userIds = append(userIds, b.TargetId)
		}

		if b.Type == "discover" {
			discoverIds = append(discoverIds, b.TargetId)
		}
	}

	print(userIds, discoverIds, err)

	return userIds, discoverIds, nil
}
