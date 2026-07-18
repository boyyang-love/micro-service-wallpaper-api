package post

import (
	"context"
	"fmt"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPostCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostCreateLogic {
	return &PostCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostCreateLogic) PostCreate(req *types.PostCreateReq) (resp *types.PostCreateRes, err error) {
	userId := fmt.Sprintf("%s", l.ctx.Value("Id"))
	if userId == "" || userId == "<nil>" {
		return &types.PostCreateRes{
			Base: types.Base{Code: 0, Msg: "请先登录"},
		}, nil
	}

	if req.Title == "" {
		return &types.PostCreateRes{
			Base: types.Base{Code: 0, Msg: "标题不能为空"},
		}, nil
	}

	if req.ImageIds == "" {
		return &types.PostCreateRes{
			Base: types.Base{Code: 0, Msg: "请至少上传一张图片"},
		}, nil
	}

	post := models.Post{
		UserId:   userId,
		Title:    req.Title,
		Content:  req.Content,
		ImageIds: req.ImageIds,
		Status:   2,
	}

	if err = l.svcCtx.DB.Create(&post).Error; err != nil {
		return nil, err
	}

	return &types.PostCreateRes{
		Base: types.Base{Code: 1, Msg: "ok"},
		Data: types.PostCreateData{Id: post.Id},
	}, nil
}
