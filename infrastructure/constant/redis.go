package constant

import "fmt"

var (
	adminLoginAuthCodeKey = "admin.login.auth.code.%s" // 后台管理登录短信验证码

	nonceKey      = "anti_replay_attack_nonce_%s"      // 防重放攻击 nonce
	idempotentKey = "anti_replay_attack_idempotent_%s" // 防幂等攻击 idempotent

	txExchangeLockKey = "tx.exchange.lock.code.%s" // 交易所交易对
)

// GetAdminLoginAuthCodeKey 获取后台管理员登录验证码
func GetAdminLoginAuthCodeKey(phone string) string {
	return fmt.Sprintf(adminLoginAuthCodeKey, phone)
}

// GetNonceKey 获取防重放攻击 nonce
func GetNonceKey(nonce string) string {
	return fmt.Sprintf(nonceKey, nonce)
}

// GetIdempotentKey 获取防幂等攻击 idempotent
func GetIdempotentKey(idempotent string) string {
	return fmt.Sprintf(idempotentKey, idempotent)
}

// GetTxExchangeLockKey 获取交易所交易对
func GetTxExchangeLockKey(code string) string {
	return fmt.Sprintf(txExchangeLockKey, code)
}
