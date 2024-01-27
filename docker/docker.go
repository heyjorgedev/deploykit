package docker

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type Connection struct {
	*client.Client
}

func NewConnection() *Connection {
	return &Connection{}
}

func (c *Connection) Open() error {
	dc, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	c.Client = dc

	return nil
}

func (c *Connection) Close() error {
	if c.Client != nil {
		if err := c.Client.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Connection) Hack() string {
	n, err := c.Client.NetworkCreate(context.Background(), "test", types.NetworkCreate{
		Attachable: true,
		Labels: map[string]string{
			"created-by": "deploykit",
		},
	})

	if err != nil {
		panic(err)
	}

	return n.ID
}

func (c *Connection) addDefaultLabels(config *container.Config) {
	if config == nil {
		config = &container.Config{}
	}

	if config.Labels == nil {
		config.Labels = map[string]string{}
	}

	config.Labels["managed-by"] = "deploykit"
}

func (c *Connection) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) (container.CreateResponse, error) {
	c.addDefaultLabels(config)
	return c.Client.ContainerCreate(ctx, config, hostConfig, networkingConfig, platform, containerName)
}
