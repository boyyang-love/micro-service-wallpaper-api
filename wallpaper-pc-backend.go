package main

import (
	"flag"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/config"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/handler"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/wallpaper-pc-backend.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(
		c.RestConf,
		rest.WithCustomCors(
			corsMidFn,
			notAllowedFn,
			[]string{
				"http://localhost:3000",
				"http://localhost:9527",
				"https://www.boyyang.cn",
				"https://boyyang.cn",
				"http://wallpaperadmin.boyyang.cn",
				"https://wallpaperadmin.boyyang.cn",
			}...,
		),
	)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

func notAllowedFn(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func corsMidFn(header http.Header) {
	header.Set("Access-Control-Allow-Headers", "X-Request-Id, Content-Type, Authorization")
	header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	header.Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
	header.Set("Access-Control-Allow-Credentials", "true")
}
