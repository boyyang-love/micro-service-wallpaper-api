// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package upload

import (
	"context"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImageUploadStreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImageUploadStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageUploadStreamLogic {
	return &ImageUploadStreamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImageUploadStreamLogic) ImageUploadStream(req *types.ImageUploadReq) (resp *types.ImageUploadRes, err error) {
	// todo: add your logic here and delete this line

	return
}
