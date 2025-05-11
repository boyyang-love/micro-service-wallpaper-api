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
	var records []types.DiscoverListInfo
	var count int64

	DB := l.svcCtx.
		DB.
		Order("created desc").
		Model(&models.Discover{})

	if req.Status != 0 {
		DB = DB.Where("status = ?", req.Status)
	} else {
		DB = DB.Where("status = ?", 2)
	}

	if req.UserId != "" {
		print(req.UserId)
		DB = DB.Or("user_id = ?", req.UserId)
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
