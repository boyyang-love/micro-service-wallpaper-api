package recommend

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecommendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecommendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecommendListLogic {
	return &RecommendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecommendListLogic) RecommendList(req *types.RecommendListReq) (resp *types.RecommendListRes, err error) {

	recommendIds, err := l.recommendId(req)
	if err != nil {
		return nil, err
	}

	uploadIds, err := l.recommend(req, recommendIds)
	if err != nil {
		return nil, err
	}

	var recommendListInfo []types.RecommendListInfo
	var count int64

	if err := l.svcCtx.
		DB.
		Order("RAND()").
		Model(&models.Upload{}).
		Select("created", "updated", "id", "file_path", "file_name").
		Where("id IN (?) and type = ?", uploadIds, req.Type).
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&recommendListInfo).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	return &types.RecommendListRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.RecommendListResData{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: recommendListInfo,
		},
	}, nil
}

func (l *RecommendListLogic) recommendId(req *types.RecommendListReq) (recommendIds []string, err error) {
	if err := l.svcCtx.
		DB.
		Model(&models.UploadRecommend{}).
		Select("recommend_id").
		Where("upload_id = ?", req.Id).
		Find(&recommendIds).
		Error; err != nil {
		return nil, err
	}

	return recommendIds, nil
}

func (l *RecommendListLogic) recommend(req *types.RecommendListReq, ids []string) (uploadIds []string, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.UploadRecommend{}).
		Select("upload_id").
		Where("recommend_id IN (?) and upload_id != ?", ids, req.Id).
		Find(&uploadIds).
		Error; err != nil {
		return nil, err
	}

	return uploadIds, nil
}
