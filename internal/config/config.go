package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type MySQLConf struct {
	Host      string
	Port      int
	Database  string
	Username  string
	Password  string
	Charset   string
	Collation string
	Timeout   string
}

type Auth struct {
	AccessSecret string
	AccessExpire int64
}

type Config struct {
	rest.RestConf
	Auth      Auth
	UploadRpc zrpc.RpcClientConf
	MySQLConf MySQLConf
}
