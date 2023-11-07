package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func GetPath(part string) string {
	savePath, _ := os.Getwd()
	savePath += string(filepath.Separator) + "picture" + string(filepath.Separator)
	switch part {
	case "head", "avatar":
		savePath += "avatar"
	case "text":
		savePath += "inset"
	case "carousel":
		savePath += "carousel"
	}
	return savePath + string(filepath.Separator)
}

func UploadPicture(c *gin.Context, part, fileName string) (err error) {
	// 解析表单，限制上传文件大小为 32MB
	err = c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		HandleError(c, http.StatusOK, "图片过大，上传失败")
		return
	}

	// 获取上传的文件
	file, handler, err := c.Request.FormFile(part)
	if err != nil {
		HandleError(c, http.StatusOK, "图片上传失败")
		return
	}
	defer file.Close()

	// 补全系统路径
	savePath := GetPath(part)
	// 判断文件夹是否存在
	_, err = os.Stat(savePath)
	if os.IsNotExist(err) {
		// 文件夹不存在
		if err = os.MkdirAll(savePath, 0755); err != nil {
			HandleError(c, http.StatusOK, "图片上传失败")
			return
		}
	}

	// 解码原始图片
	image, _, err := image.Decode(file)
	if err != nil {
		HandleError(c, http.StatusOK, "图片上传失败")
		return
	}

	// 创建并打开本地文件用于保存上传的图片
	dst, err := os.Create(savePath + fileName)
	if err != nil {
		HandleError(c, http.StatusOK, "图片上传失败")
		return
	}
	defer dst.Close()

	// 对头像图片进行压缩
	if part == "head" {
		image = resize.Resize(200, 200, image, resize.Lanczos3)
	}

	// 根据文件扩展名选择保存方式
	fileExt := strings.ToLower(path.Ext(handler.Filename))
	switch fileExt {
	case ".png":
		err = png.Encode(dst, image)
	case ".jpg", ".jpeg":
		err = jpeg.Encode(dst, image, nil)
	default:
		HandleError(c, http.StatusOK, "请上传jpg(jpeg)或png格式的图片")
		return
	}
	if err != nil {
		HandleError(c, http.StatusOK, "图片上传失败")
		return
	}
	return
}

func ShowPicture(c *gin.Context, part, fileName string) {
	savePath := GetPath(part)
	imageFile, err := os.Open(savePath + fileName)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, "图片刷新失败")
		return
	}
	defer imageFile.Close()

	imageInfo, err := imageFile.Stat()
	if err != nil {
		HandleError(c, http.StatusInternalServerError, "图片刷新失败")
		return
	}

	contentLength := imageInfo.Size()
	contentType := "image/png"

	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", contentLength))

	// 在处理响应数据之前设置正确的字符编码
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")

	_, err = io.Copy(c.Writer, imageFile)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, "头像刷新失败")
		return
	}
}
