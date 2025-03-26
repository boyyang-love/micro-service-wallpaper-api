package login

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/url"
	"strings"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignInByQqUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignInByQqUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignInByQqUrlLogic {
	return &SignInByQqUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignInByQqUrlLogic) SignInByQqUrl() (resp *types.SignInByQqUrlRes, err error) {

	uid := uuid.NewV1()
	state := strings.Replace(uid.String(), "-", "", -1)

	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", l.svcCtx.Config.QqLoginConf.AppId)
	params.Add("state", state)
	params.Add("redirect_uri", l.svcCtx.Config.QqLoginConf.RedirectURI)
	baseUrl := "https:/graph.qq.com/oauth2.0/authorize"

	u := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	return &types.SignInByQqUrlRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.SignInByQqUrlResData{
			Url:   u,
			State: state,
		},
	}, nil
}
