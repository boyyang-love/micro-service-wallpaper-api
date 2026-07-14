package user

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserLogic {
	return &ListUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUserLogic) ListUser(req *types.UserListReq) (resp *types.UserListRes, err error) {
	var records []types.UserListInfo
	var count int64

	DB := l.svcCtx.DB
	if req.Username != "" {
		DB = DB.Where("username like ?", "%"+req.Username+"%")
	}
	if req.Account != "" {
		DB = DB.Where("account like ?", "%"+req.Account+"%")
	}
	if req.Role != "" {
		DB = DB.Where("role = ?", req.Role)
	}

	if err = DB.
		Order("created desc").
		Model(&models.User{}).
		Select("created", "updated", "id", "username", "account", "avatar", "role", "address", "tel", "email", "qq", "wechat", "git_hub", "motto", "cover").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&records).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	return &types.UserListRes{
		Base: types.Base{
			Code: 1,
			Msg:  "OK",
		},
		Data: types.UserListData{
			BaseRecord: types.BaseRecord{
				Total: count,
				Limit: req.Limit,
				Page:  req.Page,
			},
			Records: records,
		},
	}, nil
}
