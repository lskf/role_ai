package repos

import (
	"context"
	"errors"
	"github.com/leor-w/kid/database/redis"
	"github.com/leor-w/kid/database/repos"
	"role_ai/infrastructure/constant"
)

type ISignerRepository interface {
	repos.IRedisRepository
	CheckNonce(nonce string) error // CheckNonce 检查 nonce 是否存在，防止重放攻击
	SaveNonce(nonce string) error  // SaveNonce 保存 nonce，防止重放攻击

	CheckIdempotent(idempotent string) error  // CheckIdempotent 检查 idempotent 是否存在，防止幂等攻击
	SaveIdempotent(idempotent string) error   // SaveIdempotent 保存 idempotent，防止幂等攻击
	DeleteIdempotent(idempotent string) error // DeleteIdempotent 删除 idempotent
}

type SignerRepository struct {
	*redis.RedisRepository `inject:""`
}

func (repo *SignerRepository) Provide(context.Context) any {
	return repo
}

// CheckNonce 检查 nonce 是否存在
func (repo *SignerRepository) CheckNonce(nonce string) error {
	exists, err := repo.Exists(constant.GetNonceKey(nonce))
	if err != nil {
		return err
	}
	if exists {
		return errors.New("nonce 已存在，疑似重放攻击")
	}
	return nil
}

// SaveNonce 保存 nonce
func (repo *SignerRepository) SaveNonce(nonce string) error {
	return repo.Set(constant.GetNonceKey(nonce), "1", 60)
}

func (repo *SignerRepository) CheckIdempotent(idempotent string) error {
	exists, err := repo.Exists(constant.GetIdempotentKey(idempotent))
	if err != nil {
		return err
	}
	if exists {
		return errors.New("idempotent 已存在，疑似幂等攻击")
	}
	return nil
}

func (repo *SignerRepository) SaveIdempotent(idempotent string) error {
	return repo.Set(constant.GetIdempotentKey(idempotent), "1", 60)
}

func (repo *SignerRepository) DeleteIdempotent(idempotent string) error {
	return repo.Del(constant.GetIdempotentKey(idempotent))
}
