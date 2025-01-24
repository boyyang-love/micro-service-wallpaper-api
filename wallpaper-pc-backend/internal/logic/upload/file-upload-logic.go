package upload

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
	return &FileUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadLogic) FileUpload(req *types.FileUploadReq) (resp *types.FileUploadRes, err error) {

	return &types.FileUploadRes{
		Base: types.Base{
			Code: 1,
			Msg:  "图片上传成功",
		},
		Data: types.FileUploadResdata{
			FileName:   req.FileName,
			Path:       req.FilePath,
			OriginPath: "",
		},
	}, err
}
