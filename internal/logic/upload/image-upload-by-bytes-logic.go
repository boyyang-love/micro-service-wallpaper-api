package upload

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImageUploadByBytesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImageUploadByBytesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageUploadByBytesLogic {
	return &ImageUploadByBytesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImageUploadByBytesLogic) ImageUploadByBytes(req *types.ImageUploadReq) (resp *types.ImageUploadRes, err error) {
	// todo: add your logic here and delete this line

	return
}
