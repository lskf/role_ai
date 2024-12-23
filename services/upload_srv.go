package services

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/leor-w/kid"
	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/errors"
	"github.com/leor-w/kid/plugin/qiniu"
	"google.golang.org/api/option"
	"io"
	"log"
	"os"
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
