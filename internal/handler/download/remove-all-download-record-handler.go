package download

import (
	"net/http"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/logic/download"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RemoveAllDownloadRecordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := download.NewRemoveAllDownloadRecordLogic(r.Context(), svcCtx)
		resp, err := l.RemoveAllDownloadRecord()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
