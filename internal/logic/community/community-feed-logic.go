// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package community

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommunityFeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommunityFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommunityFeedLogic {
	return &CommunityFeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommunityFeedLogic) CommunityFeed(req *types.CommunityFeedReq) (resp *types.CommunityFeedRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 20
	}

	var uploads []models.Upload
	var count int64

	db := l.svcCtx.DB.Model(&models.Upload{})

	// 只查已审核通过的
	db = db.Where("status = ?", 1)

	// 按类型筛选
	if req.Type != "" {
		db = db.Where("type = ?", req.Type)
	}

	// 排序
	switch req.Sort {
	case "hot":
		// 热门：按 download + view 综合排序
		db = db.Order("(download + view) desc")
	default:
		db = db.Order("created desc")
	}

	if err = db.
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&uploads).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	records := make([]types.CommunityFeedInfo, 0, len(uploads))
	for _, u := range uploads {
		// 查询点赞数
		var likeCount int64
		l.svcCtx.DB.Model(&models.Like{}).
			Where("upload_id = ? AND status = ?", u.Id, true).
			Count(&likeCount)

		// 查询用户信息
		var user models.User
		userInfo := types.CommunityUserInfo{}
		if err := l.svcCtx.DB.Model(&models.User{}).
			Where("id = ?", u.UserId).
			Select("id", "username", "avatar").
			First(&user).Error; err == nil {
			userInfo.Id = user.Id
			userInfo.Username = user.Username
			userInfo.Avatar = user.Avatar
		}

		records = append(records, types.CommunityFeedInfo{
			BaseTime: types.BaseTime{
				Created: u.Created,
				Updated: u.Updated,
			},
			Id:        u.Id,
			FileName:  u.FileName,
			FilePath:  u.FilePath,
			W:         u.W,
			H:         u.H,
			Type:      u.Type,
			View:      u.View,
			Download:  u.Download,
			LikeCount: likeCount,
			UserInfo:  userInfo,
		})
	}

	return &types.CommunityFeedRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.CommunityFeedData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: records,
		},
	}, nil
}
