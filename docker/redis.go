package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/heyjorgedev/deploykit"
)

var _ deploykit.RedisManagerService = &RedisManagerService{}

type RedisManagerService struct {
	docker *Connection
}

func NewRedisManagerService(docker *Connection) *RedisManagerService {
	return &RedisManagerService{
		docker: docker,
	}
}

func (r RedisManagerService) Create(redis *deploykit.Redis) error {
	resp, err := r.docker.ContainerCreate(
		context.TODO(),
		&container.Config{
			Image: "redis:alpine",
			Tty:   false,
		},
		nil,
		nil,
		nil,
		fmt.Sprintf("redis-%s", redis.Name),
	)
	if err != nil {
		return err
	}

	err = r.docker.ContainerStart(context.TODO(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return nil
}
