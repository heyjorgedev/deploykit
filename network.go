package deploykit

import "context"

type Network struct {
	Name              string `json:"name"`
	InternalNetworkId string `json:"-"`
}

type NetworkService interface {
	FindAll(ctx context.Context) ([]*Network, error)
}
