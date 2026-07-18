package user

import (
	"context"
	"fmt"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserStatsLogic {
	return &GetUserStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserStatsLogic) GetUserStats(req *types.UserStatsReq) (resp *types.UserStatsRes, err error) {
	userId := fmt.Sprintf("%s", l.ctx.Value("Id"))
	if userId == "" || userId == "<nil>" {
		return &types.UserStatsRes{
			Base: types.Base{Code: 0, Msg: "请先登录"},
		}, nil
	}

	// 帖子数（已发布 status=1）
	var postCount int64
	l.svcCtx.DB.Model(&models.Post{}).
		Where("user_id = ? AND status = ?", userId, 1).
		Count(&postCount)

	// 获赞数（用户所有帖子收到的点赞）
	var likeReceivedCount int64
	l.svcCtx.DB.Model(&models.Like{}).
		Where("upload_id IN (?) AND status = ?",
			l.svcCtx.DB.Model(&models.Post{}).Select("id").Where("user_id = ?", userId),
			true,
		).Count(&likeReceivedCount)

	// 评论数（用户发表的所有评论）
	var commentCount int64
	l.svcCtx.DB.Model(&models.Comment{}).
		Where("user_id = ?", userId).
		Count(&commentCount)

	// 下载数（用户帖子的下载次数总和）
	var downloadCount int64
	l.svcCtx.DB.Model(&models.Upload{}).
		Where("user_id = ? AND status = ?", userId, 1).
		Select("COALESCE(SUM(download), 0)").
		Scan(&downloadCount)

	return &types.UserStatsRes{
		Base: types.Base{Code: 1, Msg: "ok"},
		Data: types.UserStatsData{
			PostCount:         postCount,
			LikeReceivedCount: likeReceivedCount,
			CommentCount:      commentCount,
			DownloadCount:     downloadCount,
		},
	}, nil
}
