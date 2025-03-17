package login

import (
	"context"
	"errors"
	"github.com/boyyang-love/micro-service-wallpaper-api/helper"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/zeromicro/go-zero/core/collection"
	"time"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignUpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignUpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignUpLogic {
	return &SignUpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignUpLogic) SignUp(req *types.SignUpReq) (resp *types.SignUpRes, err error) {

	cache, _ := collection.NewCache(time.Minute * 5)
	code, is := cache.Get(req.Account)

	if !is {
		return nil, errors.New("不存在该账号的验证码")
	}

	if code != req.Code {
		return nil, errors.New("验证码错误")
	}

	password, err := helper.MakeHash(req.Password)
	if err != nil {
		return nil, err
	}

	if err := l.svcCtx.DB.Model(&models.User{}).Create(&models.User{
		Username: req.Username,
		Account:  req.Account,
		Password: password,
		Email:    req.Account,
	}).Error; err != nil {
		return nil, err
	}

	return &types.SignUpRes{
		Base: types.Base{
			Code: 1,
			Msg:  "账号注册成功",
		},
	}, nil
}
