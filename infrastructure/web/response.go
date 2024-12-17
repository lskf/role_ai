package web

import (
	"github.com/leor-w/kid/errors"
	"github.com/leor-w/kid/logger"

	"role_ai/infrastructure/ecode"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(data ...interface{}) *Response {
	if len(data) > 0 {
		return responseFromStatus(ecode.OK, data[0])
	}
	return responseFromStatus(ecode.OK, nil)
}

func Error(err error, data ...interface{}) *Response {
	logger.Errorf("请求处理错误: %s", err.Error())
	if len(data) > 0 {
		return responseFromStatus(errors.GetStatus(err), data[0])
	}
	return responseFromStatus(errors.GetStatus(err), nil)
}

func ParamsErr(errs ...error) *Response {
	var err error
	if len(errs) > 0 {
		err = errs[0]
		logger.Errorf("request params error: %s", err.Error())
	}
	return responseFromStatus(ecode.ReqParamInvalidErr, nil)
}

func ParamsMissingErr(errs ...error) *Response {
	var err error
	if len(errs) > 0 {
		err = errs[0]
		logger.Errorf("请求参数缺失: %s", err.Error())
	}
	return responseFromStatus(ecode.ReqParamMissErr, nil)
}

func Unauthorized(errs ...error) *Response {
	var err error
	if len(errs) > 0 {
		err = errs[0]
		logger.Errorf("权限未验证: %s", err.Error())
	}
	return responseFromStatus(ecode.NotFoundOperatorErr, err)
}

func NotFoundRoute() *Response {
	return responseFromStatus(ecode.NotFoundRouteErr, nil)
}

func ErrorWithStatus(status *errors.Status, errs ...error) *Response {
	if len(errs) > 0 {
		logger.Error(errs[0].Error())
	}
	return responseFromStatus(status, nil)
}

func responseFromStatus(status *errors.Status, data interface{}) *Response {
	return response(status.Code, status.Message, data)
}

func response(code int, msg string, data interface{}) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
