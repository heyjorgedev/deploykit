package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Connection struct {
	docker *client.Client
}

func NewConnection() *Connection {
	return &Connection{}
}

func (c *Connection) Open() error {
	dc, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	c.docker = dc

	return nil
}

func (c *Connection) Close() error {
	if c.docker != nil {
		if err := c.docker.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Connection) Hack() string {
	n, err := c.docker.NetworkCreate(context.Background(), "test", types.NetworkCreate{
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
