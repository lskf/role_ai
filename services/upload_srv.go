package services

import (
	"cloud.google.com/go/storage"
	"context"
	errors2 "errors"
	"fmt"
	"github.com/leor-w/kid"
	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/errors"
	"github.com/leor-w/kid/plugin/qiniu"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"role_ai/common"
	"role_ai/infrastructure/ecode"
	"time"
)

type UploadService struct {
	qiniu *qiniu.Qiniu `inject:""`
}

func (srv *UploadService) Provide(context.Context) any {
	return srv
}

func (srv *UploadService) UploadFormFile(ctx *kid.Context, uid int64) (string, error) {
	formFile, header, err := ctx.Request.FormFile("upload_file")
	if err != nil {
		return "", errors.New(ecode.ReqParamInvalidErr, err)
	}
	defer formFile.Close()

	//创建保存文件
	destFilePath := "./temp/" + header.Filename
	destFile, err := os.Create(destFilePath)
	if err != nil {
		return "", errors.New(ecode.InternalErr, err)
	}
	if _, err = os.Stat(destFilePath); !os.IsNotExist(err) {
		defer os.Remove(destFilePath)
	}
	defer destFile.Close()

	//复制文件
	_, err = io.Copy(destFile, formFile)
	if err != nil {
		return "", errors.New(ecode.InternalErr, err)
	}

	runMode := config.GetString("app.runMode")
	envUrl := ""
	switch runMode {
	case common.RunModeRelease:
		envUrl = ""
	default:
		envUrl = runMode
	}
	filePath := fmt.Sprintf("%s/%d/%s/%s", envUrl, uid, time.Now().Format(common.TimeFormatToDate), header.Filename)

	//将文件上传到云bucket上
	path, err := srv.UploadFileToGCS(filePath, destFilePath)
	if err != nil {
		return "", errors.New(ecode.InternalErr, err)
	}
	return path, nil
}

func (srv *UploadService) UploadUrlFile(uid int64, urlStr string, filePathType int64) (string, error) {
	// 解析 URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", errors.New(ecode.ReqParamInvalidErr, err)
	}
	// 提取 URL 中的路径部分
	filePath := parsedURL.Path
	// 获取文件名（包括扩展名）
	fileName := filepath.Base(filePath)
	// 获取文件扩展名（后缀）
	fileExt := filepath.Ext(filePath)
	if fileExt == "" {
		// 获取查询参数
		queryParams := parsedURL.Query()
		if v := queryParams["filename"]; len(v) > 0 {
			fileName = queryParams["filename"][0]
		} else {
			return "", errors.New(ecode.ReqParamInvalidErr)
		}
	}

	// 发送 GET 请求
	resp, err := http.Get(urlStr)
	if err != nil {
		return "", errors.New(ecode.InternalErr, err)
	}
	defer resp.Body.Close()
	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(ecode.InternalErr, errors2.New("获取文件失败，fileUrl:"+urlStr))
	}

	// 创建目标文件
	now := time.Now()
	localPath := fmt.Sprintf("./temp/%d%d%d/%d/%s", now.Year(), now.Month(), now.Day(), uid, fileName)
	// 提取文件目录路径
	dir := filepath.Dir(localPath)
	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 如果不存在，创建目录
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return "", errors.New(ecode.InternalErr, err)
		}
	}
	outFile, err := os.Create(localPath)
	if err != nil {
		return "", errors.New(ecode.InternalErr, err)
	}
	if _, err = os.Stat(localPath); !os.IsNotExist(err) {
		//结束时删除文件
		defer os.Remove(localPath)
	}
	defer outFile.Close()

	// 将图片内容写入文件
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return "", errors.New(ecode.InternalErr, err)
	}

	runMode := config.GetString("app.runMode")
	var envUrl, fileTypePath string

	switch runMode {
	case common.RunModeRelease:
		envUrl = ""
	default:
		envUrl = runMode
	}
	switch filePathType {
	case common.UploadFileTypeUserAvatar:
		fileTypePath = common.UploadFilePathUserAvatar
	case common.UploadFileTypeRoleAvatar:
		fileTypePath = common.UploadFilePathRoleAvatar
	case common.UploadFileTypeChatPicture:
		fileTypePath = common.UploadFilePathChatPicture
	default:
		fileTypePath = "/default"
	}
	destFile := fmt.Sprintf("%s/%s/%d/%s/%s", envUrl, time.Now().Format(common.TimeFormatToDate), uid, fileTypePath, fileName)

	//将文件上传到云bucket上
	path, err := srv.UploadFileToGCS(destFile, localPath)
	if err != nil {
		return "", errors.New(ecode.InternalErr, err)
	}
	return path, nil
}

func (srv *UploadService) UploadFileToGCS(destFile, formFile string) (string, error) {
	accountJsonUrl := config.GetString("google.storage.accountFile")
	domain := config.GetString("google.storage.domain")
	bucketName := config.GetString("google.storage.bucket")

	ctx := context.Background()
	// 使用服务账户密钥文件
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(accountJsonUrl))
	if err != nil {
		return "", errors.New(ecode.InternalErr, err)
	}
	defer client.Close()

	// 打开本地文件
	file, err := os.Open(formFile)
	if err != nil {
		log.Fatalf("打开本地文件失败: %v", err)
	}
	defer file.Close()

	obj := client.Bucket(bucketName).Object(destFile)

	// 创建上传的 Writer
	writer := obj.NewWriter(ctx)
	defer writer.Close()
	// 将文件内容写入存储桶
	if _, err = io.Copy(writer, file); err != nil {
		return "", errors.New(ecode.InternalErr, err)
	}

	return domain + "/" + bucketName + "/" + destFile, nil
}
