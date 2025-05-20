package discover

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/jinzhu/copier"
	"strings"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDiscoverListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDiscoverListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDiscoverListLogic {
	return &UserDiscoverListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDiscoverListLogic) UserDiscoverList(req *types.UserDiscoverListReq) (resp *types.UserDiscoverListRes, err error) {
	userid := fmt.Sprintf("%s", l.ctx.Value("Id"))
	var discover []models.Discover
	var records = make([]types.UserDiscoverListInfo, 0)
	var count int64
	DB := l.svcCtx.
		DB.
		Order("created desc").
		Model(&models.Discover{})

	if req.Status != 0 {
		DB = DB.Where("status = ?", req.Status)
	}

	if err = DB.
		Where("user_id = ?", userid).
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

		records = append(records, types.UserDiscoverListInfo{
			BaseTime: types.BaseTime{
				Created: d.Created,
				Updated: d.Updated,
			},
			Id:       d.Id,
			Title:    d.Title,
			Subtitle: d.Subtitle,
			Status:   d.Status,
			Images:   images,
			UserInfo: *userInfo,
		})
	}

	return &types.UserDiscoverListRes{
		Base: types.Base{},
		Data: types.UserDiscoverListData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: records,
		},
	}, nil
}

func (l *UserDiscoverListLogic) GetImages(ids []string) (images []types.UserDiscoverListImages, err error) {
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

func (l *UserDiscoverListLogic) GetUserInfo(userId string) (userInfo *types.UserDiscoverUserInfo, err error) {
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
