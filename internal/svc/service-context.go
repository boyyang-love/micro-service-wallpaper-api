package svc

import (
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-api/helper"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/config"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/uploadclient"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/user/userclient"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config        config.Config
	DB            *gorm.DB
	UploadService uploadclient.Upload
	UserService   userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := helper.ConMySQL(c.MySQLConf)
	if err != nil {
		fmt.Printf("数据库连接失败(%s)\n", err.Error())
	} else {
		fmt.Println("数据库连接成功")
	}

	uploadClient := zrpc.MustNewClient(c.UploadRpc)
	userClient := zrpc.MustNewClient(c.UserRpc)

	return &ServiceContext{
		Config:        c,
		DB:            db,
		UploadService: uploadclient.NewUpload(uploadClient),
		UserService:   userclient.NewUser(userClient),
	}
}
