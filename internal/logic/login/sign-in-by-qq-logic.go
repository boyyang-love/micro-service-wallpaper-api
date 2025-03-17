package login

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-api/helper"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignInByQqLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}
type AccessTokenInfo struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type OpenIdInfo struct {
	OpenId   string `json:"openid"`
	ClientId string `json:"client_Id"`
}

type QQUserInfo struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"figureurl_qq_1"`
}

func NewSignInByQqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignInByQqLogic {
	return &SignInByQqLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignInByQqLogic) SignInByQq(req *types.SignInByQqReq) (resp *types.SignInByQqRes, err error) {

	tokenInfo, err := l.AccessToken(req)
	if err != nil {
		return nil, err
	}

	openidInfo, err := l.OpenId(tokenInfo.AccessToken)
	if err != nil {
		return nil, err
	}

	userInfo, err := l.UserInfo(tokenInfo.AccessToken, openidInfo.OpenId)
	if err != nil {
		return nil, err
	}

	info, err := l.CreateOrUpdate(openidInfo.OpenId, userInfo)
	if err != nil {
		return nil, err
	}

	token, err := helper.NewToken(
		&helper.JwtStruct{
			Id:               info.Id,
			Username:         info.Username,
			Role:             info.Role,
			RegisteredClaims: jwt.RegisteredClaims{},
		},
		l.svcCtx.Config.Auth.AccessSecret,
		l.svcCtx.Config.Auth.AccessExpire,
	)
	if err != nil {
		return nil, err
	}

	return &types.SignInByQqRes{
		Base: types.Base{
			Code: 1,
			Msg:  "登录成功",
		},
		Data: types.SignInByQqResData{
			Token: token,
			UserInfo: types.SignInByQqUserInfo{
				Id:       info.Id,
				Username: info.Username,
				Avatar:   info.Avatar,
			},
		},
	}, nil

}

func (l *SignInByQqLogic) AccessToken(req *types.SignInByQqReq) (accessToken *AccessTokenInfo, err error) {

	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("client_id", l.svcCtx.Config.QqLoginConf.AppId)
	params.Add("client_secret", l.svcCtx.Config.QqLoginConf.AppKey)
	params.Add("code", req.Code)
	params.Add("fmt", "json")
	str := fmt.Sprintf("%s&redirect_uri=%s", params.Encode(), l.svcCtx.Config.QqLoginConf.RedirectURI)
	loginURL := fmt.Sprintf("%s?%s", "https://graph.qq.com/oauth2.0/token", str)

	response, err := http.Get(loginURL)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	var accessTokenInfo AccessTokenInfo
	err = json.Unmarshal(body, &accessTokenInfo)
	if err != nil {
		return nil, err
	}

	return &accessTokenInfo, nil

}

func (l *SignInByQqLogic) OpenId(accessToken string) (openIdInfo *OpenIdInfo, err error) {
	resp, err := http.Get(fmt.Sprintf("%s?access_token=%s&fmt=json", "https://graph.qq.com/oauth2.0/me", accessToken))
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &openIdInfo)
	if err != nil {
		return nil, err
	}

	return openIdInfo, nil
}

func (l *SignInByQqLogic) UserInfo(accessToken string, openid string) (userInfo *QQUserInfo, err error) {

	params := url.Values{}
	params.Add("access_token", accessToken)
	params.Add("openid", openid)
	params.Add("oauth_consumer_key", l.svcCtx.Config.QqLoginConf.AppKey)

	uri := fmt.Sprintf("https://graph.qq.com/user/get_user_info?%s", params.Encode())
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (l *SignInByQqLogic) CreateOrUpdate(openId string, userInfo *QQUserInfo) (info *models.User, err error) {

	user := models.User{
		Account:  openId,
		Avatar:   userInfo.Avatar,
		Username: userInfo.Nickname,
		Role:     "user",
	}

	if err = l.svcCtx.
		DB.
		Model(&models.User{}).
		Select("id", "account", "username", "avatar").
		Where("account = ?", openId).
		First(&user).
		Error; err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			if err := l.svcCtx.
				DB.
				Model(&models.User{}).
				Create(&user).
				Error; err != nil {
				return nil, err
			}
		}
	}

	if err = l.svcCtx.
		DB.
		Model(&models.User{}).
		Select("username", "avatar").
		Updates(user).
		Error; err != nil {
		return nil, err
	}

	return info, err
}
