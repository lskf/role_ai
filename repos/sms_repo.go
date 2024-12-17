package repos

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/leor-w/kid/database/redis"
	"github.com/leor-w/kid/database/repos"
	"github.com/leor-w/kid/plugin/lock"
	"role_ai/models"
	"time"
)

type ISmsRepository interface {
	repos.IRedisRepository
	LockKey(phone string, codeType int64) (ok bool, err error)
	UnLockKey(phone string, codeType int64) (err error)
	SaveLoginCode(phone, code string, codeType int64) (err error)
	GetLoginCode(phone string, codeType int64) (code string, err error)
}

type SmsRepository struct {
	*redis.RedisRepository `inject:""`
	lock                   lock.Lock `inject:""`
}

func (repo *SmsRepository) Provide(context.Context) any {
	return &SmsRepository{}
}

var (
	loginCodeLockClientKey = "login.code.lock.client.%s"
	loginCodeClientKey     = "login.code.client.%s"
)

// LockKey
// @Description: 锁定用户手机号，防止用户频繁发送短信
// @receiver repo
// @param phone
// @param codeType 1:用户
// @return ok
// @return err
func (repo *SmsRepository) LockKey(phone string, codeType int64) (ok bool, err error) {
	switch codeType {
	case 1:
		ok, err = repo.lock.Lock(fmt.Sprintf(loginCodeLockClientKey, phone), time.Minute*1)
	default:
		ok, err = false, nil
	}
	return
}

// UnLockKey
// @Description: 解除用户手机号锁定
// @receiver repo
// @param phone
// @param codeType
// @return err
func (repo *SmsRepository) UnLockKey(phone string, codeType int64) (err error) {
	switch codeType {
	case 1:
		err = repo.lock.Unlock(fmt.Sprintf(loginCodeLockClientKey, phone))
	}
	return
}

// SaveLoginCode
// @Description: 保存验证码
// @receiver repo
// @param phone
// @param code
// @param codeType
// @return err
func (repo *SmsRepository) SaveLoginCode(phone, code string, codeType int64) (err error) {
	verifyCode := models.VerifyCode{
		Code:    code,
		Issued:  time.Now().Unix(),
		Expired: time.Now().Add(time.Minute * 5).Unix(),
	}
	verifyCodeStr, _ := json.Marshal(verifyCode)
	switch codeType {
	case 1:
		err = repo.RDB.Set(fmt.Sprintf(loginCodeClientKey, phone), verifyCodeStr, time.Minute*5).Err()
	}
	return
}

// GetLoginCode
// @Description: 获取验证码
// @receiver repo
// @param phone
// @param codeType
// @return code
// @return err
func (repo *SmsRepository) GetLoginCode(phone string, codeType int64) (code string, err error) {
	switch codeType {
	case 1:
		code, err = repo.RDB.Get(fmt.Sprintf(loginCodeClientKey, phone)).Result()
	}
	return
}
