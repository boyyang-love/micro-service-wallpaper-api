package upload

import (
	"context"
	"fmt"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserImageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserImageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserImageListLogic {
	return &UserImageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserImageListLogic) UserImageList(req *types.UserImageListReq) (resp *types.UserImageListRes, err error) {
	userId := fmt.Sprintf("%s", l.ctx.Value("Id"))
	if userId == "" || userId == "<nil>" {
		return &types.UserImageListRes{
			Base: types.Base{Code: 0, Msg: "未登录"},
		}, nil
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 20
	}

	var uploadInfo []Upload
	var imageInfo []types.ImageInfo
	var count int64

	DB := l.svcCtx.
		DB.
		Preload("Tags").
		Preload("Category").
		Preload("Recommend").
		Preload("Group").
		Preload("Album").
		Model(&Upload{}).
		Where("user_id = ?", userId).
		Order("created desc")

	if req.Type != "" {
		DB = DB.Where("type = ?", req.Type)
	}

	if err := DB.
		Count(&count).
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Find(&uploadInfo).
		Offset(-1).
		Error; err != nil {
		return nil, err
	}

	_ = copier.Copy(&imageInfo, &uploadInfo)

	return &types.UserImageListRes{
		Base: types.Base{Code: 1, Msg: "ok"},
		Data: types.ImageInfoResdata{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: imageInfo,
		},
	}, nil
}
