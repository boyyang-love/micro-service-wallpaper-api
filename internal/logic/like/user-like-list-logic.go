package like

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type Like struct {
	models.Like
	Upload Upload `gorm:"foreignkey:upload_id;references:id"`
}

type Upload struct {
	models.Upload
}

func (l *Like) TableName() string {
	return "like"
}

func (l *Upload) TableName() string {
	return "upload"
}

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
	var like []Like
	var records = make([]types.UserLikeListRecord, 0)
	var count int64

	DB := l.svcCtx.
		DB.
		Order("updated desc").
		Preload("Upload", func(db *gorm.DB) *gorm.DB {
			db = db.Select("id", "file_path", "file_name", "w", "h", "type")
			//if req.Type != "" {
			//	db = db.Where("type = ?", req.Type)
			//}
			return db
		}).
		Model(&Like{})

	if req.Type != "" {
		DB = DB.Where("type = ?", req.Type)
	}

	if err = DB.
		Where("user_id = ? and status = ?", userId, true).
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&like).
		Offset(-1).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	for _, v := range like {
		records = append(records, types.UserLikeListRecord{
			BaseTime: types.BaseTime{
				Created: v.Created,
				Updated: v.Updated,
			},
			Id:       v.Id,
			FileId:   v.Upload.Id,
			FilePath: v.Upload.FilePath,
			FileName: v.Upload.FileName,
			W:        v.Upload.W,
			H:        v.Upload.H,
			Type:     v.Upload.Type,
		})
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
