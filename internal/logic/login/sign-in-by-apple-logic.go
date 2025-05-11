package login

import (
	"context"
	"errors"
	"github.com/boyyang-love/micro-service-wallpaper-api/helper"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignInByAppleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignInByAppleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignInByAppleLogic {
	return &SignInByAppleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignInByAppleLogic) SignInByApple(req *types.SignInByAppleReq) (resp *types.SignInByAppleRes, err error) {
	err, info := l.IsExist(req.AppleId)
	if err != nil {
		return nil, err
	}

	if info == nil {
		var user = &models.User{
			Username: req.Name,
			AppleId:  req.AppleId,
			Email:    req.Email,
			Role:     "user",
		}
		res := l.svcCtx.DB.Model(&models.User{}).Create(&user)

		if res.Error != nil {
			return nil, res.Error
		}

		token, err := helper.NewToken(&helper.JwtStruct{
			Id:               user.Id,
			Username:         req.Name,
			Role:             "user",
			RegisteredClaims: jwt.RegisteredClaims{},
		}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)

		if err != nil {
			return nil, err
		}

		return &types.SignInByAppleRes{
			Base: types.Base{
				Code: 1,
				Msg:  "登录成功",
			},
			Data: types.SignInByAppleResData{
				Token: token,
				UserInfo: types.SignInByAppleUserInfo{
					Id:       user.Id,
					Username: req.Name,
					Avatar:   user.Avatar,
				},
			},
		}, nil
	}

	token, err := helper.NewToken(&helper.JwtStruct{
		Id:               info.Id,
		Username:         info.Username,
		Role:             info.Role,
		RegisteredClaims: jwt.RegisteredClaims{},
	}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)

	if err != nil {
		return nil, err
	}

	return &types.SignInByAppleRes{
		Base: types.Base{
			Code: 1,
			Msg:  "登录成功",
		},
		Data: types.SignInByAppleResData{
			Token: token,
			UserInfo: types.SignInByAppleUserInfo{
				Id:       info.Id,
				Username: info.Username,
				Avatar:   info.Avatar,
			},
		},
	}, nil
}

func (l *SignInByAppleLogic) IsExist(appleId string) (err error, info *models.User) {
	var user models.User
	if err = l.svcCtx.
		DB.
		Model(&models.User{}).
		Select("id", "username", "apple_id", "avatar", "role").
		Where("apple_id = ?", appleId).
		First(&user).
		Error; err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return err, nil
		}
	}

	return nil, &user
}
