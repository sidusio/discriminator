package settings

import (
	"context"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	AppName       = "discriminator"
	ReverseDomain = "io.sidus"
)

const (
	templatesPath      = "templates-path"
	templatesExtension = "templates-extension"

	containerLabel           = "container-label"
	includeStoppedContainers = "include-stopped-containers"

	runInterval = "run-interval"
)

type Settings struct {
	v *viper.Viper
}

func NewSettings(_ context.Context) (Settings, error) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(
		"-", "_",
	))
	v.SetEnvPrefix(AppName)

	setDefaults(v)

	v.AutomaticEnv()
	return Settings{v: v}, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault(templatesPath, "/templates")
	v.SetDefault(templatesExtension, ".tmpl")

	v.SetDefault(containerLabel, ReverseDomain+"."+AppName)
	v.SetDefault(includeStoppedContainers, false)

	v.SetDefault(runInterval, 5*time.Minute)
}

func (s Settings) TemplatesPath() string {
	return s.v.GetString(templatesPath)
}

func (s Settings) TemplatesExtension() string {
	return s.v.GetString(templatesExtension)
}

func (s Settings) ContainerLabel() string {
	return s.v.GetString(containerLabel)
}

func (s Settings) IncludeStoppedContainers() bool {
	return s.v.GetBool(includeStoppedContainers)
}

func (s Settings) RunInterval() time.Duration {
	return s.v.GetDuration(runInterval)
}
