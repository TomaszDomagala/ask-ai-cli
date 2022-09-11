package cmd

import (
	"context"

	"github.com/spf13/viper"
)

const GlobalConfigKey = "globalConfig"

func GetGlobalConfig(ctx context.Context) *viper.Viper {
	return ctx.Value(GlobalConfigKey).(*viper.Viper)
}
func SetGlobalConfig(ctx context.Context, conf *viper.Viper) context.Context {
	return context.WithValue(ctx, GlobalConfigKey, conf)
}

const ProviderConfigKey = "providerConfig"

func GetProviderConfig(ctx context.Context) *viper.Viper {
	return ctx.Value(ProviderConfigKey).(*viper.Viper)
}
func SetProviderConfig(ctx context.Context, conf *viper.Viper) context.Context {
	return context.WithValue(ctx, ProviderConfigKey, conf)
}
