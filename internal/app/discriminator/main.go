package discriminator

import (
	"context"

	"github.com/pkg/errors"

	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"

	"sidus.io/discriminator/internal/pkg/docker"
	"sidus.io/discriminator/internal/pkg/parsing"
	"sidus.io/discriminator/internal/pkg/templates"
)

var (
	applicationLabel = "io.sidus.discriminator"
	templatesPath    = "./templates"
	extension        = ".tmpl"
	includeStopped   = false
)

// Run runs the application for one iteration
func Run() error {
	ctx := context.Background()
	logrus.WithContext(ctx).Infof("Building templates directory...")
	tmpls, err := templates.LoadTemplatesFromPath(ctx, templatesPath, extension)
	if err != nil {
		return errors.Wrapf(err, "failed to load templates")
	}
	templateDirectory, err := templates.NewDirectory(ctx, tmpls, extension)
	if err != nil {
		return err
	}
	logrus.WithContext(ctx).Infof("Built templates directory with %d templates", templateDirectory.Count(ctx))

	parser, err := parsing.NewParser(ctx, templateDirectory)
	if err != nil {
		return errors.Wrapf(err, "failed to create parser")
	}

	logrus.WithContext(ctx).Infof("Connecting to docker client")
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return errors.Wrapf(err, "failed to create docker client from environment")
	}
	dockerService, err := docker.NewService(ctx, dockerClient)
	if err != nil {
		return err
	}
	containers, err := dockerService.GetContainers(ctx, includeStopped)
	if err != nil {
		return err
	}
	logrus.WithContext(ctx).Infof("Retrieved %d containers from the docker client", len(containers))
	for _, container := range containers {
		value, ok := container.Labels[applicationLabel]
		if ok {
			logrus.WithContext(ctx).Infof("Processing container %s (%s) with options: %v", container.Name, container.ID, value)
			modifiers, err := parser.Process(ctx, value, templates.ContainerData{
				Labels: container.Labels,
			})
			if err != nil {
				logrus.WithError(err).Errorf("encountered error while processing container %s (%s)", container.Name, container.ID)
			}
			newLabels := modifiers.Apply(container.Labels)

			if !stringMapEquals(newLabels, container.Labels) {
				logrus.WithContext(ctx).Infof("Updating %s (%s) with new labels", container.Name, container.ID)
				err = dockerService.SetLabels(ctx, container.ID, newLabels)
				if err != nil {
					logrus.WithError(err).Errorf("encountered error while setting labels on container %s (%s)", container.Name, container.ID)
				}
			}
		}
	}
	return nil
}

// stringMapEquals compare two maps of type map[string]string
func stringMapEquals(a, b map[string]string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for key, value := range a {
		if bValue, ok := b[key]; bValue != value || !ok {
			return false
		}
	}
	return true
}
