package docker

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

var timeout = 30 * time.Second

type Service struct {
	dockerClient Client
}

// NewService creates a dervice to be used for docker communication
func NewService(_ context.Context, dockerClient Client) (*Service, error) {
	c := Service{
		dockerClient: dockerClient,
	}
	return &c, nil
}

// GetContainers retrieves a list of containers from the configured docker endpoint
func (s *Service) GetContainers(ctx context.Context, includeStopped bool) ([]Container, error) {
	dockerContainers, err := s.dockerClient.ContainerList(ctx, types.ContainerListOptions{
		All: includeStopped,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list containers")
	}
	containers := make([]Container, len(dockerContainers))
	for i, dockerContainer := range dockerContainers {
		containers[i] = Container{
			Name:   firstOrEmpty(dockerContainer.Names),
			ID:     dockerContainer.ID,
			Labels: dockerContainer.Labels,
		}
	}
	return containers, nil
}

// SetLabels removes the old container and creat a new, identical one with the specified labels.
//
// New container will not have the same id
func (s *Service) SetLabels(ctx context.Context, containerID string, labels map[string]string) error {
	ctx = context.WithValue(ctx, "containerID", containerID)
	ctx = context.WithValue(ctx, "newContainerLabels", labels)

	logrus.WithContext(ctx).Debugf("Inspecting container %s", containerID)
	container, err := s.dockerClient.ContainerInspect(ctx, containerID)
	if err != nil {
		return errors.Wrapf(err, "inspection failed for container with id: %s", containerID)
	}
	ctx = context.WithValue(ctx, "oldContainerLabels", container.Config.Labels)
	ctx = context.WithValue(ctx, "containerName", container.Name)

	logrus.WithContext(ctx).Debugf("Stopping container %s", containerID)
	err = s.dockerClient.ContainerStop(ctx, containerID, &timeout)
	if err != nil {
		// TODO: should maybe be handled? what happens on timeout?
		return errors.Wrapf(err, "failed to stop container %s", containerID)
	}

	newName := container.Name + "-old"
	logrus.WithContext(ctx).Debugf("Changing name of container %s from %s to %s", containerID, container.Name, newName)
	err = s.dockerClient.ContainerRename(ctx, containerID, newName)
	if err != nil {
		// TODO: retry mechanism
		return errors.Wrapf(err, "failed to rename container %s from %s to %s", containerID, container.Name, newName)
	}

	logrus.WithContext(ctx).Debugf("creating new container with name: %s", container.Name)
	config := container.Config
	// Setting labels
	config.Labels = labels
	netConfig := network.NetworkingConfig{
		EndpointsConfig: container.NetworkSettings.Networks,
	}
	newID, err := s.dockerClient.ContainerCreate(ctx, config, container.HostConfig, &netConfig, container.Name)
	if err != nil {
		// TODO: restore old container
		return errors.Wrapf(
			err,
			"failed to create new container with name: %s, old container can be restored from (id: %s, name:%s)",
			container.Name, containerID, newName,
		)
	}

	if container.State.Running {
		err = s.dockerClient.ContainerStart(ctx, newID.ID, types.ContainerStartOptions{})
		if err != nil {
			// TODO: restore old container
			return errors.Wrapf(
				err,
				"failed to start new container with (name: %s, id: %s), old container can be restored from (id: %s, name:%s)",
				container.Name, newID, containerID, newName,
			)
		}
	}

	logrus.WithContext(ctx).Debugf("Removing old container with name: %s and id: %s", container.Name, containerID)
	err = s.dockerClient.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		RemoveVolumes: false,
		RemoveLinks:   false,
		Force:         false,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to remove old container (%s) with name: %s", containerID, newName)
	}
	return nil
}

// Utility functions

func firstOrEmpty(ss []string) string {
	if len(ss) > 0 {
		return ss[0]
	} else {
		return ""
	}
}
