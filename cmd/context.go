package cmd

import (
	"context"
	"github.com/spf13/afero"

	"github.com/spf13/viper"
)

const GlobalConfigKey = "globalConfig"

func GetGlobalConfig(ctx context.Context) *viper.Viper {
	return ctx.Value(GlobalConfigKey).(*viper.Viper)
}
func SetGlobalConfig(ctx context.Context, conf *viper.Viper) context.Context {
	return context.WithValue(ctx, GlobalConfigKey, conf)
}

const FileSystemKey = "fileSystem"

func GetFs(ctx context.Context) afero.Fs {
	return ctx.Value(FileSystemKey).(afero.Fs)
}

func SetFs(ctx context.Context, fs afero.Fs) context.Context {
	return context.WithValue(ctx, FileSystemKey, fs)
}
