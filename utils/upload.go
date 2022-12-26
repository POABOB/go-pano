package utils

import (
	"fmt"
	"mime/multipart"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 使用協程上傳
func SaveWithGo(ctx *gin.Context, base_dir string, files []*multipart.FileHeader) error {
	startTime := time.Now().UnixMicro()
	// 處理檔案的channel
	working := make(chan *multipart.FileHeader)
	// 處理錯誤的channel
	failures := make(chan error)

	// 定義waiting Group 避免程式直接跑完
	var wg sync.WaitGroup

	// goroutine worker * 5 分工上傳，然後把結果寫出來
	for i := 0; i < 5; i++ {
		wg.Add(1)
		// worker 處理檔案
		go worker(&wg, working, failures, ctx, base_dir)
	}

	// goroutine 將單個file傳給chan，讓協程自動執行
	for _, file := range files {
		working <- file
	}

	// 關閉 chan
	close(working)
	// 等待 goroutine 執行完畢
	wg.Wait()
	waitingTime := time.Now().UnixMicro() - startTime
	fmt.Println("Upload taking times in", float64(waitingTime)/1000.00, "ms...")

	// 檢查看看有沒有錯誤
	close(failures)
	if err := <-failures; err != nil {
		return err
	}
	return nil
}

// worker上傳檔案，並記錄錯誤
func worker(wg *sync.WaitGroup, working <-chan *multipart.FileHeader, failures chan<- error, ctx *gin.Context, base_dir string) {
	// 最後執行 Done
	defer wg.Done()
	for f := range working {
		err := ctx.SaveUploadedFile(f, base_dir+"/"+f.Filename)
		if err != nil {
			fmt.Println(err.Error())
			failures <- err
		}
	}
}
