package like

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDownloadAndLikeSummaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDownloadAndLikeSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDownloadAndLikeSummaryLogic {
	return &UserDownloadAndLikeSummaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDownloadAndLikeSummaryLogic) UserDownloadAndLikeSummary() (resp *types.UserDownloadAndLikeSummaryRes, err error) {
	var userId = fmt.Sprintf("%s", l.ctx.Value("Id"))

	downloadCount, err := l.DownloadSummary(userId)
	if err != nil {
		return nil, err
	}

	likeCount, err := l.LikeSummary(userId)
	if err != nil {
		return nil, err
	}

	discoverCount, err := l.DiscoverSummary(userId)
	if err != nil {
		return nil, err
	}

	blockCount, err := l.BlockSummary(userId)
	if err != nil {
		return nil, err
	}

	return &types.UserDownloadAndLikeSummaryRes{
		Base: types.Base{
			Code: 1,
			Msg:  "ok",
		},
		Data: types.UserDownloadAndLikeSummaryData{
			Download: downloadCount,
			Like:     likeCount,
			Discover: discoverCount,
			Block:    blockCount,
		},
	}, nil
}

func (l *UserDownloadAndLikeSummaryLogic) DownloadSummary(userId string) (count int64, err error) {
	var ids []string
	if err = l.svcCtx.
		DB.
		Model(&models.Download{}).
		Where("user_id = ? and type != ?", userId, "AVATAR").
		Select("download_id").
		Find(&ids).
		Error; err != nil {
		return 0, err
	}

	if err = l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Where("id in (?)", ids).
		Count(&count).
		Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (l *UserDownloadAndLikeSummaryLogic) LikeSummary(userId string) (count int64, err error) {
	var ids []string
	if err = l.svcCtx.
		DB.
		Model(&models.Like{}).
		Where("user_id = ? and status = ? and type != ?", userId, true, "AVATAR").
		Select("upload_id").
		Find(&ids).
		Error; err != nil {
		return 0, err
	}

	if err = l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Where("id in (?)", ids).
		Count(&count).
		Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (l *UserDownloadAndLikeSummaryLogic) DiscoverSummary(userId string) (count int64, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.Discover{}).
		Where("user_id = ?", userId).
		Count(&count).
		Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (l *UserDownloadAndLikeSummaryLogic) BlockSummary(userId string) (count int64, err error) {
	if err = l.svcCtx.
		DB.
		Model(&models.Block{}).
		Where("user_id = ?", userId).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
