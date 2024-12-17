package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/leor-w/kid"
	"github.com/leor-w/kid/middleware"
	"net/http"
	kid_mid "role_ai/infrastructure/middleware"
	"time"
)

func InitRouter(kid *kid.Kid) {
	kid.Use(middleware.CrosWithConfig(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{http.MethodPost, http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodHead, http.MethodOptions, http.MethodPatch},
		AllowHeaders:    []string{"Authorization", "Content-Type", "Content-Length", "X-CSRF-Token", "Token", "sign", "X-Custom-Header", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
		MaxAge:          12 * time.Hour,
	}))
	kid.NoRoute(kid_mid.NoRoute())
	api := kid.Group("/api")
	v1 := api.Group("/v1")
	InitAppRouter(v1)
}
