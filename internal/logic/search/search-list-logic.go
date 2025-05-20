package search

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchListLogic {
	return &SearchListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchListLogic) SearchList(req *types.SearchListReq) (resp *types.SearchListRes, err error) {
	var count int64
	var uploads []models.Upload
	records := make([]types.SearchListInfo, 0)
	ids, err := l.GetSearchTagIds(req.Keywords)
	if err != nil {
		return nil, err
	}

	if err = l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Select("id", "created", "updated", "file_name", "file_path", "w", "h", "download", "view").
		Where("file_name LIKE ? and type = ?", "%"+req.Keywords+"%", req.Type).
		Or("id in (?) and type = ? ", ids, req.Type).
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&uploads).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	err = copier.Copy(&records, &uploads)
	if err != nil {
		return nil, err
	}

	return &types.SearchListRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.SearchListData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: records,
		},
	}, nil
}

func (l *SearchListLogic) GetSearchTagIds(keywords string) (ids []string, err error) {
	var tagIds []string
	if err := l.svcCtx.
		DB.
		Model(&models.Tag{}).
		Select("id").
		Where("name LIKE ?", "%"+keywords+"%").
		Find(&tagIds).
		Error; err != nil {
		return nil, err
	}

	if len(tagIds) > 0 {
		if err := l.svcCtx.
			DB.
			Model(&models.UploadTag{}).
			Select("upload_id").
			Where("tag_id in (?)", tagIds).
			Find(&ids).
			Error; err != nil {
			return nil, err
		}
	}

	return ids, nil
}
