package upload

import (
	"context"
	"os"
	"path/filepath"
	"strconv"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChunkCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChunkCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChunkCheckLogic {
	return &ChunkCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChunkCheckLogic) ChunkCheck(req *types.ChunkCheckReq) (resp *types.ChunkCheckRes, err error) {
	hashDir := filepath.Join(chunkDir, req.FileHash)

	uploadedList := make([]int, 0)
	for i := 0; i < req.TotalChunks; i++ {
		chunkPath := filepath.Join(hashDir, strconv.Itoa(i))
		if _, err := os.Stat(chunkPath); err == nil {
			uploadedList = append(uploadedList, i)
		}
	}

	logx.Infof("检查分片状态: 文件哈希=%s, 总分片=%d, 已上传=%d", req.FileHash, req.TotalChunks, len(uploadedList))

	return &types.ChunkCheckRes{
		Base: types.Base{Code: 1, Msg: "ok"},
		Data: types.ChunkCheckResData{
			FileHash:       req.FileHash,
			TotalChunks:    req.TotalChunks,
			UploadedChunks: uploadedList,
			UploadedCount:  len(uploadedList),
		},
	}, nil
}
