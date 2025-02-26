package upload

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/uploadclient"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImageDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImageDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageDeleteLogic {
	return &ImageDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImageDeleteLogic) ImageDelete(req *types.ImageDeleteReq) (resp *types.ImageDeleteRes, err error) {
	if err := l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Where("id = ?", req.Id).
		Delete(&models.Upload{}).
		Error; err != nil {
		return nil, err
	}

	_, err = l.svcCtx.UploadService.ImageDelete(l.ctx, &uploadclient.ImageDeleteReq{
		BucketName: req.BucketName,
		Paths:      req.Paths,
	})
	if err != nil {
		return nil, err
	}

	return &types.ImageDeleteRes{
		Base: types.Base{
			Code: 1,
			Msg:  "删除成功",
		},
	}, nil
}
