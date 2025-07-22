package block

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateBlockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateBlockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBlockLogic {
	return &CreateBlockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateBlockLogic) CreateBlock(req *types.BlockCreateReq) (resp *types.BlockCreateRes, err error) {
	userid := fmt.Sprintf("%s", l.ctx.Value("Id"))
	if err = l.svcCtx.
		DB.
		Debug().
		Model(&models.Block{}).
		Create(&models.Block{
			TargetId: req.TargetId,
			Type:     req.Type,
			UserId:   userid,
		}).Error; err != nil {
		return nil, err
	}

	return &types.BlockCreateRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
	}, nil
}
