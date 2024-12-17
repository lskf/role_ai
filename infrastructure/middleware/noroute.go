package middleware

import (
	"github.com/leor-w/kid"

	"role_ai/infrastructure/web"
)

func NoRoute() kid.HandleFunc {
	return func(ctx *kid.Context) interface{} {
		return web.NotFoundRoute()
	}
}
