package like

import (
	"net/http"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/logic/like"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RemoveAllLikeRecordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := like.NewRemoveAllLikeRecordLogic(r.Context(), svcCtx)
		resp, err := l.RemoveAllLikeRecord()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
