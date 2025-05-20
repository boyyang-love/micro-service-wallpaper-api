package like

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
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
	var records = make([]types.UserLikeListRecord, 0)
	var count int64

	uploadIds, err := l.UploadIds(userId)
	if err != nil {
		return nil, err
	}

	DB := l.svcCtx.
		DB.
		Order(fmt.Sprintf("'%s'", strings.Join(uploadIds, "','"))).
		Model(&models.Upload{})

	if req.Type != "" {
		DB = DB.Where("type = ?", req.Type)
	}

	if err = DB.
		Where("id in (?)", uploadIds).
		Select("created", "updated", "id", "file_path", "file_name", "w", "h", "type").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&records).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	var sortedRecords []types.UserLikeListRecord

	for _, uploadId := range uploadIds {
		for _, record := range records {
			if uploadId == record.Id {
				sortedRecords = append(sortedRecords, record)
				break
			}
		}
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
			Records: sortedRecords,
		},
	}, nil
}

func (l *UserLikeListLogic) UploadIds(userId string) (ids []string, err error) {

	if err = l.svcCtx.
		DB.
		Order("updated desc").
		Model(&models.Like{}).
		Select("upload_id").
		Where("user_id = ? and status = ?", userId, true).
		Find(&ids).
		Error; err != nil {
		return ids, err
	}

	return ids, nil
}
