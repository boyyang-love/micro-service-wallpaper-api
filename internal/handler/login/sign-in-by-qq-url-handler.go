package login

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"strings"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/logic/login"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SignInByQqUrlHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := login.NewSignInByQqUrlLogic(r.Context(), svcCtx)
		resp, err := l.SignInByQqUrl()
		uid := uuid.NewV1()
		state := strings.Replace(uid.String(), "-", "", -1)

		params := url.Values{}
		params.Add("response_type", "code")
		params.Add("client_id", svcCtx.Config.QqLoginConf.AppId)
		params.Add("state", state)
		params.Add("redirect_uri", svcCtx.Config.QqLoginConf.RedirectURI)
		baseUrl := "https:/graph.qq.com/oauth2.0/authorize"
		u := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

		http.Redirect(w, r, u, http.StatusFound)

		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
