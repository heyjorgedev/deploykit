package sqlite

type AppService struct{}

func NewAppService(db *DB) *AppService {
	return &AppService{}
}
