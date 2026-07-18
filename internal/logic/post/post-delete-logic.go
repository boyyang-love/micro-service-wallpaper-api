package post

import (
	"context"
	"fmt"
	"strings"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-models/models"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/uploadclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPostDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostDeleteLogic {
	return &PostDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostDeleteLogic) PostDelete(req *types.PostDeleteReq) (resp *types.PostDeleteRes, err error) {
	userId := fmt.Sprintf("%s", l.ctx.Value("Id"))
	if userId == "" || userId == "<nil>" {
		return &types.PostDeleteRes{
			Base: types.Base{Code: 0, Msg: "请先登录"},
		}, nil
	}

	// 查找帖子（校验所有权）
	var post models.Post
	if err = l.svcCtx.DB.Model(&models.Post{}).
		Where("id = ? AND user_id = ?", req.Id, userId).
		First(&post).Error; err != nil {
		return &types.PostDeleteRes{
			Base: types.Base{Code: 0, Msg: "帖子不存在或无权删除"},
		}, nil
	}

	// 解析图片 ID
	var imageIds []string
	if post.ImageIds != "" {
		imageIds = strings.Split(post.ImageIds, ",")
	}

	// 查询图片信息（用于删除存储文件）
	type filePaths struct {
		OriginFilePath string
		FilePath       string
	}
	var paths []filePaths
	if len(imageIds) > 0 {
		l.svcCtx.DB.Model(&models.Upload{}).
			Where("id IN (?)", imageIds).
			Select("origin_file_path", "file_path").
			Find(&paths)
	}

	// 开始事务删除
	tx := l.svcCtx.DB.Begin()

	// 1. 删除帖子
	if err = tx.Where("id = ?", post.Id).Delete(&models.Post{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 2. 删除帖子的评论
	if err = tx.Where("post_id = ?", post.Id).Delete(&models.Comment{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 3. 删除帖子的点赞
	if err = tx.Where("upload_id = ?", post.Id).Delete(&models.Like{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 4. 删除图片关联数据（标签、分类、推荐）
	if len(imageIds) > 0 {
		if err = tx.Where("upload_id IN (?)", imageIds).Delete(&models.UploadTag{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err = tx.Where("upload_id IN (?)", imageIds).Delete(&models.UploadCategory{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err = tx.Where("upload_id IN (?)", imageIds).Delete(&models.UploadRecommend{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// 5. 删除图片记录
		if err = tx.Where("id IN (?)", imageIds).Delete(&models.Upload{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	// 6. 删除存储文件（事务外执行，不影响数据一致性）
	if len(paths) > 0 {
		var allPaths []string
		for _, p := range paths {
			if p.OriginFilePath != "" {
				allPaths = append(allPaths, p.OriginFilePath)
			}
			if p.FilePath != "" {
				allPaths = append(allPaths, p.FilePath)
			}
		}
		if len(allPaths) > 0 {
			_, _ = l.svcCtx.UploadService.CosDelete(l.ctx, &uploadclient.ImageDeleteReq{
				BucketName: "wallpaper",
				Paths:      allPaths,
			})
		}
	}

	return &types.PostDeleteRes{
		Base: types.Base{Code: 1, Msg: "ok"},
	}, nil
}
