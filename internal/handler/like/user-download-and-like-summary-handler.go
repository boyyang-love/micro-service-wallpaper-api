package like

import (
	"net/http"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/logic/like"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserDownloadAndLikeSummaryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := like.NewUserDownloadAndLikeSummaryLogic(r.Context(), svcCtx)
		resp, err := l.UserDownloadAndLikeSummary()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
