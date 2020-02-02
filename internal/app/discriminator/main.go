package discriminator

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"

	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"

	"sidus.io/discriminator/internal/pkg/docker"
	"sidus.io/discriminator/internal/pkg/parsing"
	"sidus.io/discriminator/internal/pkg/settings"
	"sidus.io/discriminator/internal/pkg/templates"
)

func Start() error {
	ctx := context.WithValue(context.Background(), "phase", "setup")

	logrus.WithContext(ctx).Infof("Loading settings")
	s, err := settings.NewSettings(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to load settings")
	}
	logrus.WithContext(ctx).Infof("Settings loaded")

	logrus.WithContext(ctx).Infof("Setting up necessary services")
	dockerService, parser, err := setup(ctx, s)
	if err != nil {
		return errors.Wrapf(err, "failed during setup")
	}
	defer func() {
		err := dockerService.Close()
		if err != nil {
			logrus.WithError(err).Errorf("Could not close docker service")
		}
	}()
	logrus.WithContext(ctx).Infof("Setup completed")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer signal.Stop(c)

	ctx = context.WithValue(ctx, "phase", "operating")
	stop := false
	for {
		logrus.WithContext(ctx).Infof("Starting iteration...")
		ctx := context.WithValue(ctx, "runStartedAt", time.Now())
		err := run(ctx, dockerService, parser, s)
		if err != nil {
			return err
		}
		logrus.WithContext(ctx).Infof("Iteration completed, sleeping for %.0f minutes.", s.RunInterval().Minutes())

		select {
		case <-c:
			stop = true
			logrus.WithContext(ctx).Infof("Received stop signal")
			break
		case <-time.After(s.RunInterval()):
			break
		}

		if stop {
			break
		}
	}
	return nil
}

// Creates all services needed to run the application
func setup(ctx context.Context, s settings.Settings) (*docker.Service, parsing.Parser, error) {
	logrus.WithContext(ctx).Infof("Building templates directory...")
	tmpls, err := templates.LoadTemplatesFromPath(ctx, s.TemplatesPath(), s.TemplatesExtension())
	if err != nil {
		return nil, parsing.Parser{}, errors.Wrapf(err, "failed to load templates")
	}
	templateDirectory, err := templates.NewDirectory(ctx, tmpls, s.TemplatesExtension())
	if err != nil {
		return nil, parsing.Parser{}, errors.Wrapf(err, "failed to create template directory")
	}
	logrus.WithContext(ctx).Infof("Built templates directory with %d templates", templateDirectory.Count(ctx))

	parser, err := parsing.NewParser(ctx, templateDirectory)
	if err != nil {
		return nil, parsing.Parser{}, errors.Wrapf(err, "failed to create parser")
	}

	logrus.WithContext(ctx).Infof("Connecting to docker client")
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return nil, parsing.Parser{}, errors.Wrapf(err, "failed to create docker client from environment")
	}
	dockerService, err := docker.NewService(ctx, dockerClient)
	if err != nil {
		return nil, parsing.Parser{}, errors.Wrapf(err, "failed to create docker service")
	}
	return dockerService, parser, nil
}

// Run runs the application for one iteration
func run(ctx context.Context, dockerService *docker.Service, parser parsing.Parser, s settings.Settings) error {
	containers, err := dockerService.GetContainers(ctx, s.IncludeStoppedContainers())
	if err != nil {
		return err
	}
	logrus.WithContext(ctx).Infof("Retrieved %d containers from the docker client", len(containers))

	for _, container := range containers {
		value, ok := container.Labels[s.ContainerLabel()]
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
					logrus.WithError(err).Errorf(
						"encountered error while setting labels on container %s (%s)",
						container.Name,
						container.ID,
					)
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
