package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/heyjorgedev/deploykit"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type MachineService struct {
	docker *client.Client
}

func NewMachineService() *MachineService {
	return &MachineService{}
}

func (s *MachineService) Open() (err error) {
	s.docker, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	return nil
}

func (s *MachineService) Close() error {
	if s.docker != nil {
		return s.docker.Close()
	}

	return nil
}

func (s *MachineService) CreateMachine(ctx context.Context, w io.Writer, machine *deploykit.Machine) error {
	if err := ensureImageExists(ctx, w, s.docker, machine.Image); err != nil {
		return err
	}

	return createContainer(ctx, w, s.docker, machine)
}

func ensureImageExists(ctx context.Context, w io.Writer, docker *client.Client, imageName string) error {
	fmt.Println("Ensuring image exists:", imageName)

	_, _, err := docker.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		if client.IsErrNotFound(err) {
			return pullImage(ctx, w, docker, imageName)
		}
		return err
	}

	return nil
}

func pullImage(ctx context.Context, w io.Writer, docker *client.Client, imageName string) error {
	fmt.Println("Pulling image:", imageName)
	r, err := docker.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return err
	}
	defer r.Close()

	io.Copy(w, r)

	return nil
}

func createContainer(ctx context.Context, w io.Writer, docker *client.Client, machine *deploykit.Machine) error {
	fmt.Println("Creating container:", machine.Image)

	resp, err := docker.ContainerCreate(
		ctx,
		&container.Config{
			Image: machine.Image,
		},
		&container.HostConfig{
			AutoRemove: machine.AutoDestroy,
			RestartPolicy: container.RestartPolicy{
				Name: container.RestartPolicyUnlessStopped,
			},
		},
		&network.NetworkingConfig{},
		&v1.Platform{},
		machine.ID,
	)
	if err != nil {
		return err
	}

	fmt.Println("Starting container:", machine.ID)
	if err := docker.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}

	return nil
}
