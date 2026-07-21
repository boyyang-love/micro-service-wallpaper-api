package upload

import (
	"sync"
)

var (
	// 全局上传信号量
	uploadSemaphore chan struct{}
	// 初始化锁
	initOnce sync.Once
)

// InitUploadLimiter 初始化上传限流器
func InitUploadLimiter(maxConcurrent int) {
	initOnce.Do(func() {
		if maxConcurrent <= 0 {
			maxConcurrent = 10
		}
		uploadSemaphore = make(chan struct{}, maxConcurrent)
	})
}

// AcquireUpload 获取上传许可
func AcquireUpload() {
	if uploadSemaphore != nil {
		uploadSemaphore <- struct{}{}
	}
}

// ReleaseUpload 释放上传许可
func ReleaseUpload() {
	if uploadSemaphore != nil {
		<-uploadSemaphore
	}
}
