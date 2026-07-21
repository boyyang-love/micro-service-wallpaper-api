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

type QqLoginConf struct {
	AppId       string
	AppKey      string
	RedirectURI string
}

type WechatLoginConf struct {
	AppId  string
	Secret string
}

type CorsOriginConf []string

type UploadConf struct {
	MaxFileSize   int64 `json:",default=10485760"` // 10MB 默认限制
	MaxConcurrent int   `json:",default=10"`       // 最大并发数
}

type Config struct {
	rest.RestConf
	Auth            Auth
	MySQLConf       MySQLConf
	QqLoginConf     QqLoginConf
	WechatLoginConf WechatLoginConf
	CorsOriginConf  CorsOriginConf
	UploadRpc       zrpc.RpcClientConf
	UserRpc         zrpc.RpcClientConf
	UploadConf      UploadConf
}
