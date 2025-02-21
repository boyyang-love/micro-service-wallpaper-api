package main

import (
	"flag"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/config"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/handler"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/wallpaper-pc-backend.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
