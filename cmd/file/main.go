package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"im/config"
	"im/pkg/logger"
	"im/pkg/util"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	// 初始化配置
	config.Init("config.yaml")

	// 初始化日志
	logger.Init(config.GetFileServer().LogFilePath, config.GetFileServer().LogTarget, config.GetFileServer().LogLevel)

	router := gin.Default()
	router.Static("/file", config.GetFileServer().ResourcePath)

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// single file
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusOK, Response{
				Code:    1001,
				Message: err.Error(),
			})
			return
		}

		filenames := strings.Split(file.Filename, ".")
		name := strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + util.RandString(30) + "." + filenames[len(filenames)-1]
		filePath := config.GetFileServer().ResourcePath + name
		if err = c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusOK, Response{
				Code:    1001,
				Message: err.Error(),
			})
			return
		}

		url := fmt.Sprintf("http://%s/file/%s", config.GetFileServer().WideAddr, filePath)
		c.JSON(http.StatusOK, Response{
			Code:    0,
			Message: "success",
			Data:    map[string]string{"url": url},
		})
	})
	_ = router.Run(config.GetFileServer().LocalAddr)
}
