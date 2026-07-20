package login

import (
	"net/http"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/logic/login"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SignInByQqUrlHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := login.NewSignInByQqUrlLogic(r.Context(), svcCtx)
		resp, err := l.SignInByQqUrl()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
