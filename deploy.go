package deploykit

type Redis struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RedisManagerService interface {
	Create(redis *Redis) error
}
