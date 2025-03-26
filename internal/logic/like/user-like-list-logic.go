package like

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLikeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLikeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLikeListLogic {
	return &UserLikeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLikeListLogic) UserLikeList(req *types.UserLikeListReq) (resp *types.UserLikeListRes, err error) {
	var userId = fmt.Sprintf("%s", l.ctx.Value("Id"))
	var records []types.UserLikeListRecord

	uploadIds, count, err := l.UploadIds(req, userId)
	if err != nil {
		return nil, err
	}

	if err = l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Where("id in (?)", uploadIds).
		Select("created", "updated", "id", "file_path", "file_name", "w", "h", "type").
		Find(&records).
		Error; err != nil {
		return nil, err
	}

	return &types.UserLikeListRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.UserLikeListData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: records,
		},
	}, nil
}

func (l *UserLikeListLogic) UploadIds(req *types.UserLikeListReq, userId string) (ids []string, count int64, err error) {
	if err = l.svcCtx.
		DB.
		Order("created desc").
		Model(&models.Like{}).
		Select("upload_id").
		Where("user_id = ? and status = ?", userId, true).
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&ids).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return ids, count, err
	}

	return ids, count, nil
}
