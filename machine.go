package deploykit

import "context"

type Machine struct {
	ID    string
	Image string

	Env         map[string]string
	AutoDestroy bool
}

type MachineService interface {
	CreateMachine(ctx context.Context, machine *Machine) error
}
