// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2

package types

type Base struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type ImageDeleteReq struct {
	BucketName string   `json:"bucket_name"`
	Id         string   `json:"id"`
	Paths      []string `json:"paths"`
}

type ImageDeleteRes struct {
	Base
}

type ImageInfo struct {
	Id             string `json:"id"`
	Hash           string `json:"hash"`
	FileName       string `json:"file_name"`
	OriginFileSize int64  `json:"origin_file_size"`
	FileSize       int64  `json:"file_size"`
	OriginType     string `json:"origin_type"`
	FileType       string `json:"file_type"`
	OriginFilePath string `json:"origin_file_path"`
	FilePath       string `json:"file_path"`
	Type           string `json:"type"`
	W              int    `json:"w"`
	H              int    `json:"h"`
	Status         int    `json:"status"`
	UserId         string `json:"user_id"`
}

type ImageInfoReq struct {
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
	FileName string `form:"file_name,optional"`
	Type     string `form:"type,optional"`
	Status   int    `form:"status,optional"`
}

type ImageInfoRes struct {
	Base
	Data ImageInfoResdata `json:"data"`
}

type ImageInfoResdata struct {
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
	Total   int64       `json:"total"`
	Records []ImageInfo `json:"records"`
}

type ImageUpdateReq struct {
	Id       string `json:"id"`
	FileName string `json:"file_name"`
	Type     string `json:"type"`
	Status   int    `json:"status"`
}

type ImageUpdateRes struct {
	Base
}

type ImageUploadReq struct {
	FileName   string `form:"file_name"`
	Type       string `form:"type,optional"`
	RootDir    string `form:"root_dir"`
	Dir        string `form:"dir"`
	BucketName string `form:"bucket_name"`
	Quality    uint32 `form:"quality"`
	Status     int    `form:"status"`
}

type ImageUploadRes struct {
	Base
	Data ImageUploadResdata `json:"data"`
}

type ImageUploadResdata struct {
	FileName   string `json:"file_name"`
	Path       string `json:"path"`
	OriginPath string `json:"origin_path"`
}

type SignInReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInRes struct {
	Base
	Data SignInResData `json:"data"`
}

type SignInResData struct {
	Token    string                `json:"token"`
	UserInfo SignInResDataUserInfo `json:"user_info"`
}

type SignInResDataUserInfo struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Motto    string `json:"motto"`
	Address  string `json:"address"`
	Tel      string `json:"tel"`
	Email    string `json:"email"`
	QQ       string `json:"qq"`
	Wechat   string `json:"wechat"`
	GitHub   string `json:"git_hub"`
	Role     string `json:"role"`
	Avatar   string `json:"avatar"`
	Cover    string `json:"cover"`
}

type SignUpReq struct {
}

type SignUpRes struct {
}
