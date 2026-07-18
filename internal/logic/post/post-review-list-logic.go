package post

import (
	"context"
	"strings"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostReviewListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPostReviewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostReviewListLogic {
	return &PostReviewListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostReviewListLogic) PostReviewList(req *types.PostReviewListReq) (resp *types.PostReviewListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 10
	}

	var posts []models.Post
	var count int64

	db := l.svcCtx.DB.Model(&models.Post{})
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}

	if err = db.
		Order("created DESC").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&posts).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	records := make([]types.PostReviewListItem, 0, len(posts))
	for _, p := range posts {
		images := l.getImages(p.ImageIds)
		userInfo := l.getUserInfo(p.UserId)

		records = append(records, types.PostReviewListItem{
			BaseTime: types.BaseTime{
				Created: p.Created,
				Updated: p.Updated,
			},
			Id:       p.Id,
			Title:    p.Title,
			Content:  p.Content,
			ImageIds: p.ImageIds,
			Images:   images,
			Status:   p.Status,
			RejectReason: p.RejectReason,
			UserInfo: userInfo,
		})
	}

	return &types.PostReviewListRes{
		Base: types.Base{Code: 1, Msg: "ok"},
		Data: types.PostReviewListData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: records,
		},
	}, nil
}

func (l *PostReviewListLogic) getImages(imageIds string) []types.PostImageInfo {
	if imageIds == "" {
		return []types.PostImageInfo{}
	}
	ids := strings.Split(imageIds, ",")
	var uploads []models.Upload
	l.svcCtx.DB.Model(&models.Upload{}).
		Select("id", "file_name", "file_path", "w", "h").
		Where("id IN (?)", ids).
		Find(&uploads)

	images := make([]types.PostImageInfo, 0, len(uploads))
	for _, u := range uploads {
		images = append(images, types.PostImageInfo{
			Id:       u.Id,
			FilePath: u.FilePath,
			FileName: u.FileName,
			W:        u.W,
			H:        u.H,
		})
	}
	return images
}

func (l *PostReviewListLogic) getUserInfo(userId string) types.PostUserInfo {
	var user models.User
	if err := l.svcCtx.DB.Model(&models.User{}).
		Where("id = ?", userId).
		Select("id", "username", "avatar").
		First(&user).Error; err != nil {
		return types.PostUserInfo{}
	}
	return types.PostUserInfo{
		Id:       user.Id,
		Username: user.Username,
		Avatar:   user.Avatar,
	}
}
