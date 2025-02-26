// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2

package handler

import (
	"net/http"
	"time"

	login "github.com/boyyang-love/micro-service-wallpaper-api/internal/handler/login"
	upload "github.com/boyyang-love/micro-service-wallpaper-api/internal/handler/upload"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/signin",
				Handler: login.SignInHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/signup",
				Handler: login.SignUpHandler(serverCtx),
			},
		},
		rest.WithTimeout(20000*time.Millisecond),
		rest.WithMaxBytes(20971520),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/image/delete",
				Handler: upload.ImageDeleteHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/image/info",
				Handler: upload.ImageInfoHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/image/update",
				Handler: upload.ImageUpdateHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/image/upload",
				Handler: upload.ImageUploadHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithTimeout(20000*time.Millisecond),
		rest.WithMaxBytes(20971520),
	)
}
