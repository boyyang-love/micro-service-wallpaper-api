package svc

import (
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/config"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/uploadclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	UploadService uploadclient.Upload
}

func NewServiceContext(c config.Config) *ServiceContext {
	uploadClient := zrpc.MustNewClient(c.UploadRpc)
	return &ServiceContext{
		Config:        c,
		UploadService: uploadclient.NewUpload(uploadClient),
	}
}
