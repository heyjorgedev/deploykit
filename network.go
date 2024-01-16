package deploykit

type Network struct {
	ID                uint16 `json:"id"`
	Name              string `json:"name"`
	InternalNetworkId string `json:"-"`
}

type NetworkEntityManager interface {
}
