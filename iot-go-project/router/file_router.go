/*
Copyright 2024 - 2025 Zen HuiFer

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type FileApi struct{}

func (*FileApi) UpdateFile(c *gin.Context) {
	// 绑定上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		// 如果有错误，返回错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	zap.S().Infof(file.Filename)
	uuidStr := uuid.New().String()
	// 定义文件保存的路径
	dst := "./fileupdate/" + uuidStr + "_" + file.Filename

	// 打开上传的文件
	fileReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot open file."})
		return
	}
	defer func(fileReader multipart.File) {
		err := fileReader.Close()
		if err != nil {
			zap.S().Error(err.Error())
		}
	}(fileReader)

	// 创建保存的文件并打开
	out, err := os.Create(dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create file."})
		return
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			zap.S().Error(err.Error())
		}
	}(out)

	// 将上传的文件内容写入到新文件中
	if _, err = io.Copy(out, fileReader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot write file."})
		return
	}

	// 返回文件的保存路径
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file_path": dst})
}

// DownloadFile 根据提供的文件路径下载文件
func (*FileApi) DownloadFile(c *gin.Context) {

	// 获取请求中的文件路径参数
	filePath := c.Query("path")

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting file info"})
		return
	}

	// 获取文件名
	fileName := fileInfo.Name()

	// 设置HTTP响应头
	c.Header("Content-Disposition", "attachment; filename="+strings.ReplaceAll(fileName, "\"", ""))
	c.Header("Content-Type", "application/octet-stream")

	c.File(filePath)
}
