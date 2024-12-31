package services

import (
	"context"
	"encoding/json"
	"github.com/leor-w/kid/errors"
	"github.com/leor-w/kid/plugin/smscode"
	"role_ai/infrastructure/ecode"
	"role_ai/infrastructure/tools"
	"role_ai/models"
	"role_ai/repos"
	"strconv"
	"time"
)

type SmsService struct {
	smsRepo repos.ISmsRepository `inject:""`
	sms     *smscode.Ali         `inject:""`
}

func (srv *SmsService) Provide(_ context.Context) interface{} {
	return srv
}

// SendCode
// @Description: 发送短信验证码
// @receiver srv
// @param phone
// @return error
func (srv *SmsService) SendCode(phone string) error {
	code := tools.RandomInt64InRange(100000, 999999)
	codeStr := strconv.FormatInt(code, 10)
	//反正频繁操作
	ok, err := srv.smsRepo.LockKey(phone, 1)
	if err != nil {
		return errors.New(ecode.SendSMSCodeIntervalErr, err)
	}
	if !ok {
		return errors.New(ecode.SendSMSCodeIntervalErr)
	}
	//保存
	err = srv.smsRepo.SaveLoginCode(phone, codeStr, 1)
	if err != nil {
		return errors.New(ecode.DatabaseErr, err)
	}
	//发送
	if err = srv.sms.SendSMS(phone, codeStr); err != nil {
		return errors.New(ecode.SMSSendErr, err)
	}
	return nil
}

// VerifyCode
// @Description: 校验短信验证码
// @receiver srv
// @param phone
// @param code
// @return bool
// @return error
func (srv *SmsService) VerifyCode(phone string, code string) (bool, error) {
	//获取验证码
	verifyCodeStr, err := srv.smsRepo.GetLoginCode(phone, 1)
	if err != nil {
		return false, errors.New(ecode.DatabaseErr, err)
	}
	verifyCode := models.VerifyCode{}
	err = json.Unmarshal([]byte(verifyCodeStr), &verifyCode)
	if err != nil {
		return false, errors.New(ecode.DataProcessingErr, err)
	}

	if verifyCode.Code != code || verifyCode.Expired < time.Now().Unix() {
		return false, errors.New(ecode.LoginAuthCodeErr)
	}
	//解除用户手机号码锁定
	_ = srv.smsRepo.UnLockKey(phone, 1)
	return true, nil
}
