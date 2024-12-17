package services

import (
	"context"
	"github.com/leor-w/kid/plugin/qiniu"
)

type UploadService struct {
	qiniu *qiniu.Qiniu `inject:""`
}

func (srv *UploadService) Provide(context.Context) any {
	return srv
}

func (srv *UploadService) GetQiniuUploadToken() string {
	return srv.qiniu.UploadToken()
}
