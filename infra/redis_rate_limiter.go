package infra

import (
	"fmt"
	"go-dms/config"
)

func ResetLogin(ip string) error {
	key := fmt.Sprintf("login:rl:%s", ip)
	return config.Client.Del(config.Ctx, key).Err()
}

func ResetRefreshToken(ip string, refreshToken string) error {
	key := fmt.Sprintf("refresh:rl:%s:%s", ip, refreshToken)
	return config.Client.Del(config.Ctx, key).Err()
}
