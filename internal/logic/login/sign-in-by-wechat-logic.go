package login

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/boyyang-love/micro-service-wallpaper-api/helper"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type SignInByWechatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type WechatSession struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid,omitempty"`
	Errcode    int    `json:"errcode,omitempty"`
	Errmsg     string `json:"errmsg,omitempty"`
}

func NewSignInByWechatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignInByWechatLogic {
	return &SignInByWechatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignInByWechatLogic) SignInByWechat(req *types.SignInByWechatReq) (resp *types.SignInByWechatRes, err error) {

	session, err := l.Jscode2session(req.Code)
	if err != nil {
		return nil, err
	}

	if err = l.CreateOrUpdate(session.Openid); err != nil {
		return nil, err
	}

	info, token, err := l.InfoAndToken(session.Openid)
	if err != nil {
		return nil, err
	}

	return &types.SignInByWechatRes{
		Base: types.Base{
			Code: 1,
			Msg:  "登录成功",
		},
		Data: types.SignInByWechatResData{
			Token: token,
			UserInfo: types.SignInByWechatUserInfo{
				Id:       info.Id,
				Username: info.Username,
				Avatar:   info.Avatar,
			},
		},
	}, nil

}

func (l *SignInByWechatLogic) Jscode2session(code string) (*WechatSession, error) {
	params := url.Values{}
	params.Add("appid", l.svcCtx.Config.WechatLoginConf.AppId)
	params.Add("secret", l.svcCtx.Config.WechatLoginConf.Secret)
	params.Add("js_code", code)
	params.Add("grant_type", "authorization_code")

	uri := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?%s", params.Encode())
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var session WechatSession
	if err = json.Unmarshal(body, &session); err != nil {
		return nil, err
	}

	if session.Errcode != 0 {
		return nil, errors.New(session.Errmsg)
	}

	return &session, nil
}

func (l *SignInByWechatLogic) CreateOrUpdate(openId string) error {
	var user models.User
	err := l.svcCtx.DB.
		Model(&models.User{}).
		Select("id").
		Where("open_id = ?", openId).
		First(&user).Error

	if err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			if err := l.svcCtx.DB.
				Model(&models.User{}).
				Create(&models.User{
					OpenId:   openId,
					Username: fmt.Sprintf("微信用户%s", openId[:8]),
					Role:     "user",
				}).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (l *SignInByWechatLogic) InfoAndToken(openId string) (*models.User, string, error) {
	var user models.User
	if err := l.svcCtx.DB.
		Model(&models.User{}).
		Where("open_id = ?", openId).
		Select("id", "username", "avatar", "role").
		First(&user).Error; err != nil {
		return nil, "", err
	}

	token, err := helper.NewToken(
		&helper.JwtStruct{
			Id:               user.Id,
			Username:         user.Username,
			Role:             user.Role,
			RegisteredClaims: jwt.RegisteredClaims{},
		},
		l.svcCtx.Config.Auth.AccessSecret,
		l.svcCtx.Config.Auth.AccessExpire,
	)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}
