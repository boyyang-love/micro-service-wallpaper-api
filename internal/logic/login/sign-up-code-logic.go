package login

import (
	"context"
	"errors"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/email/pb/email"
	"github.com/zeromicro/go-zero/core/collection"
	"gorm.io/gorm"
	"math/rand"
	"time"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignUpCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignUpCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignUpCodeLogic {
	return &SignUpCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignUpCodeLogic) SignUpCode(req *types.SignUpCodeReq) (resp *types.SignUpCodeRes, err error) {

	res := l.svcCtx.
		DB.
		Model(&models.User{}).
		Select("account").
		Where("account = ?", req.Account).
		First(&models.User{})

	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if res.RowsAffected == 0 {
		code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))

		emailRpcRes, err := l.svcCtx.EmailService.SendEmailCode(
			l.ctx,
			&email.SendEmailCodeReq{
				Email:   req.Account,
				Code:    code,
				Subject: "壁纸收藏家",
				Title:   "壁纸收藏家账号注册",
				Time:    5,
			},
		)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		fmt.Println(emailRpcRes.Base.Code, emailRpcRes.Base.Msg)

		cache, err := collection.NewCache(time.Minute * 5)
		if err != nil {
			return nil, err
		}

		cache.Set(req.Account, code)

		return &types.SignUpCodeRes{
			Base: types.Base{
				Code: 1,
				Msg:  "验证码发送成功",
			},
		}, nil
	}

	return nil, errors.New("该账号已存在")
}
