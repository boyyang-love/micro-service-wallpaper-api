// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2

package types

type Base struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type FileUploadReq struct {
	FileName string `json:"file_name,optional"`
	FilePath string `json:"file_path,optional"`
}

type FileUploadRes struct {
	Base
	Data FileUploadResdata `json:"data"`
}

type FileUploadResdata struct {
	FileName   string `json:"file_name"`
	Path       string `json:"path"`
	OriginPath string `json:"origin_path"`
}
