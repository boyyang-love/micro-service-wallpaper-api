package login

import (
	"context"
	"errors"
	"github.com/boyyang-love/micro-service-wallpaper-api/helper"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/user/pb/user"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignInLogic {
	return &SignInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignInLogic) SignIn(req *types.SignInReq) (resp *types.SignInRes, err error) {
	password, err := helper.MakeHash(req.Password)
	if err != nil {
		return nil, err
	}

	if l.Is(req.Username) {
		userInfo := models.User{}
		if err := l.svcCtx.
			DB.
			Model(&models.User{}).
			Select("id", "username", "password").
			Where("username = ? and password = ?", req.Username, password).
			First(&userInfo).
			Error; errors.As(err, &gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}

		info, err := l.svcCtx.UserService.UserInfo(l.ctx, &user.UserInfoReq{Id: userInfo.Id})
		if err != nil {
			return nil, err
		}

		token, err := helper.NewToken(&helper.JwtStruct{
			Id:               info.Data.Id,
			Username:         info.Data.Username,
			Role:             info.Data.Role,
			RegisteredClaims: jwt.RegisteredClaims{},
		},
			l.svcCtx.Config.Auth.AccessSecret,
			l.svcCtx.Config.Auth.AccessExpire,
		)
		if err != nil {
			return nil, errors.New("token生成失败")
		}

		return &types.SignInRes{
			Base: types.Base{
				Code: 1,
				Msg:  "登录成功",
			},
			Data: types.SignInResData{
				Token: token,
				UserInfo: types.SignInResDataUserInfo{
					Id:       info.Data.Id,
					Username: info.Data.Username,
					Motto:    info.Data.Motto,
					Address:  info.Data.Address,
					Tel:      info.Data.Tel,
					Email:    info.Data.Email,
					QQ:       info.Data.QQ,
					Wechat:   info.Data.Wechat,
					GitHub:   info.Data.GitHub,
					Role:     info.Data.Role,
					Avatar:   info.Data.Avatar,
					Cover:    info.Data.Cover,
				},
			},
		}, nil

	} else {
		return nil, errors.New("用户不存在")
	}
}

// 判断用户名是否存在
func (l *SignInLogic) Is(username string) (is bool) {
	// 定义一个User结构体变量
	userModel := models.User{}
	// 使用gorm的DB对象，查询User表中的username字段，判断用户名是否存在
	if err := l.svcCtx.
		DB.
		Model(&models.User{}).
		Select("username").
		Where("username = ?", username).
		First(&userModel).
		Error; errors.As(err, &gorm.ErrRecordNotFound) {
		// 如果查询不到，返回false
		return false
	}

	// 如果查询到，返回true
	return true
}
