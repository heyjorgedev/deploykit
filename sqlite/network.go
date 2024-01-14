package sqlite

import (
	"context"

	"github.com/heyjorgedev/deploykit"
)

type NetworkService struct {
	db *DB
}

func NewNetworkService(db *DB) *NetworkService {
	return &NetworkService{db}
}

func (s *NetworkService) FindAll(ctx context.Context) ([]*deploykit.Network, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, "SELECT name, internal_network_id FROM networks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	networks := make([]*deploykit.Network, 0)
	for rows.Next() {
		var network deploykit.Network
		if err := rows.Scan(&network.Name, &network.InternalNetworkId); err != nil {
			return nil, err
		}
		networks = append(networks, &network)
	}

	return networks, nil
}
