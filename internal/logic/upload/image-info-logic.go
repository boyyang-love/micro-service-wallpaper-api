package upload

import (
	"context"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ImageInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type Upload struct {
	models.Upload
	Tags      []Tag       `json:"tags" gorm:"many2many:upload_tag;"`
	Category  []Category  `json:"category" gorm:"many2many:upload_category;"`
	Recommend []Recommend `json:"recommend" gorm:"many2many:upload_recommend;"`
	Group     []Group     `json:"group" gorm:"many2many:upload_group;"`
}

type Tag struct {
	models.Tag
}

type Category struct {
	models.Category
}

type Recommend struct {
	models.Recommend
}

type Group struct {
	models.Group
}

func NewImageInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageInfoLogic {
	return &ImageInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImageInfoLogic) ImageInfo(req *types.ImageInfoReq) (resp *types.ImageInfoRes, err error) {
	var uploadInfo []Upload
	var imageInfo []types.ImageInfo
	var count int64
	DB := l.svcCtx.
		DB.
		Preload("Tags").
		Preload("Category").
		Preload("Recommend").
		Preload("Group").
		Model(&Upload{}).
		Order("created  desc")

	if req.Status == 1 || req.Status == 2 {
		DB = DB.Where("status = ? ", req.Status)
	}

	if req.FileName != "" {
		uploadIds, err := l.UploadIds(req)
		if err != nil {
			return nil, err
		}
		DB = DB.Where("file_name LIKE ? OR id in (?)", "%"+req.FileName+"%", uploadIds)
	}

	if req.Type != "" {
		DB = DB.Where("type = ? ", req.Type)
	}

	if err := DB.
		Count(&count).
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Find(&uploadInfo).
		Offset(-1).
		Error; err != nil {
		return nil, err
	}

	_ = copier.Copy(&imageInfo, &uploadInfo)

	return &types.ImageInfoRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.ImageInfoResdata{
			BaseRecord: types.BaseRecord{
				Page:  req.Page,
				Limit: req.Limit,
				Total: count,
			},
			Records: imageInfo,
		},
	}, nil
}

func (l *ImageInfoLogic) UploadIds(req *types.ImageInfoReq) (ids []string, err error) {

	var tagIds []string
	if err = l.svcCtx.
		DB.
		Model(&models.Tag{}).
		Select("id").
		Where("name LIKE ?", "%"+req.FileName+"%").
		Find(&tagIds).
		Error; err != nil {
		return nil, err
	}

	var uploadIds []string
	if err = l.svcCtx.
		DB.
		Model(&models.UploadTag{}).
		Select("upload_id").
		Where("tag_id in (?)", tagIds).
		Find(&uploadIds).
		Error; err != nil {
		return nil, err
	}

	return uploadIds, nil
}
