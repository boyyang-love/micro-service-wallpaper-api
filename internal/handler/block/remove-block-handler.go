package block

import (
	"net/http"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/logic/block"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RemoveBlockHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BlockRemoveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := block.NewRemoveBlockLogic(r.Context(), svcCtx)
		resp, err := l.RemoveBlock(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
