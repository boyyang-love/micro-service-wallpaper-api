package tag

import (
	"context"
	"errors"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTagLogic {
	return &CreateTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTagLogic) CreateTag(req *types.CreateTagReq) (resp *types.CreateTagRes, err error) {
	userid := fmt.Sprintf("%s", l.ctx.Value("Id"))

	is, err := l.IsTagExist(req)
	if err != nil {
		return nil, err
	}

	if is {
		return &types.CreateTagRes{
			Base: types.Base{
				Code: 2,
				Msg:  "该标签已经存在",
			},
		}, nil
	}

	if err = l.svcCtx.
		DB.
		Model(&models.Tag{}).
		Create(&models.Tag{
			Name:   req.Name,
			Type:   req.Type,
			Sort:   req.Sort,
			UserId: userid,
		}).Error; err != nil {
		return nil, err
	}

	return &types.CreateTagRes{
		Base: types.Base{
			Code: 1,
			Msg:  "标签创建成功",
		},
	}, nil
}

func (l *CreateTagLogic) IsTagExist(req *types.CreateTagReq) (is bool, err error) {
	var tag models.Tag
	if err = l.svcCtx.
		DB.
		Model(&models.Tag{}).
		Select("id", "name", "type").
		Where("name = ? and type = ?", req.Name, req.Type).
		First(&tag).
		Error; err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
