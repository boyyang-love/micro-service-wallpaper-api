package upload

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const chunkDir = "/tmp/upload-chunks"

func init() {
	os.MkdirAll(chunkDir, 0755)
}

type ChunkUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChunkUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChunkUploadLogic {
	return &ChunkUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChunkUploadLogic) ChunkUpload(req *types.ChunkUploadReq) (resp *types.ChunkUploadRes, err error) {
	r := l.ctx.Value("httpRequest").(*http.Request)
	r.Body = http.MaxBytesReader(nil, r.Body, 2<<20)

	file, _, err := r.FormFile("chunk")
	if err != nil {
		return nil, fmt.Errorf("获取分片文件失败: %v", err)
	}
	defer file.Close()

	hashDir := filepath.Join(chunkDir, req.FileHash)
	if err := os.MkdirAll(hashDir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %v", err)
	}

	chunkPath := filepath.Join(hashDir, strconv.Itoa(req.ChunkIndex))
	chunkFile, err := os.Create(chunkPath)
	if err != nil {
		return nil, fmt.Errorf("创建分片文件失败: %v", err)
	}
	defer chunkFile.Close()

	if _, err := io.Copy(chunkFile, file); err != nil {
		return nil, fmt.Errorf("保存分片失败: %v", err)
	}

	logx.Infof("分片上传成功: 文件=%s, 分片=%d/%d", req.FileName, req.ChunkIndex+1, req.TotalChunks)

	return &types.ChunkUploadRes{
		Base: types.Base{Code: 1, Msg: "分片上传成功"},
		Data: types.ChunkUploadResData{ChunkIndex: req.ChunkIndex, Uploaded: true},
	}, nil
}
