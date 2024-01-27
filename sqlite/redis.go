package sqlite

import "github.com/heyjorgedev/deploykit"

var _ deploykit.RedisManagerService = &RedisManagerService{}

type RedisManagerService struct {
	db           *DB
	orchestrator deploykit.RedisManagerService
}

func NewRedisManagerService(db *DB, orchestrator deploykit.RedisManagerService) *RedisManagerService {
	return &RedisManagerService{
		db:           db,
		orchestrator: orchestrator,
	}
}

func (r RedisManagerService) Create(redis *deploykit.Redis) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = r.orchestrator.Create(redis)
	if err != nil {
		return err
	}

	record, err := tx.Exec("INSERT INTO storage_redis (name) VALUES (?)", redis.Name)
	if err != nil {
		return err
	}

	id, err := record.LastInsertId()
	if err != nil {
		return err
	}

	redis.ID = int(id)
	return nil
}
