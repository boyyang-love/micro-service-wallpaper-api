package account

import (
	"net/http"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/logic/account"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CancelAccountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CancelAccountReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := account.NewCancelAccountLogic(r.Context(), svcCtx)
		resp, err := l.CancelAccount(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
