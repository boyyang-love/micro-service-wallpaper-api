package user

import (
	"context"
	"errors"

	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"gorm.io/gorm"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.UserInfoReq) (resp *types.UserInfoRes, err error) {
	userId, _ := l.ctx.Value("Id").(string)
	if userId == "" {
		return nil, errors.New("未登录")
	}

	var user models.User
	if err := l.svcCtx.DB.
		Model(&models.User{}).
		Where("id = ?", userId).
		First(&user).Error; err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return &types.UserInfoRes{
		Base: types.Base{Code: 1, Msg: "ok"},
		Data: types.UserInfoResData{
			Id:       user.Id,
			Username: user.Username,
			Avatar:   user.Avatar,
			Cover:    user.Cover,
			Motto:    user.Motto,
			Role:     user.Role,
		},
	}, nil
}
