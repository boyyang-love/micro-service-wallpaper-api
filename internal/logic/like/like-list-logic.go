package like

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeListLogic {
	return &LikeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeListLogic) LikeList(req *types.LikeListReq) (resp *types.LikeListRes, err error) {

	userid := fmt.Sprintf("%s", l.ctx.Value("Id"))

	var likeList []string
	if err := l.svcCtx.
		DB.
		Model(&models.Like{}).
		Select("upload_id").
		Where("user_id = ? and status = ?", userid, true).
		Find(&likeList).
		Error; err != nil {
		return nil, err
	}

	return &types.LikeListRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: likeList,
	}, nil
}
