package ctxutil

import (
	"context"

	"github.com/piplabs/story-guardian/internal/config"
)

type contentKey string

const (
	ctxContentAppConfigKey   contentKey = "appConfig"
	ctxContentAccessTokenKey contentKey = "accessToken"
)

func GetAppConfig(ctx context.Context) *config.AppConfig {
	conf, _ := ctx.Value(ctxContentAppConfigKey).(*config.AppConfig)
	return conf
}

func WithAppConfig(ctx context.Context, config *config.AppConfig) context.Context {
	return context.WithValue(ctx, ctxContentAppConfigKey, config)
}

func GetAccessToken(ctx context.Context) string {
	token, _ := ctx.Value(ctxContentAccessTokenKey).(string)
	return token
}

func WithAccessToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, ctxContentAccessTokenKey, token)
}
