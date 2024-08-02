package deploykit

import "context"

type Team struct {
	ID   int
	Name string
}

type TeamService interface {
	FindTeams(ctx context.Context, filter TeamFilter) ([]*Team, int, error)
	FindTeamByID(ctx context.Context, id int) (*Team, error)
	CreateTeam(ctx context.Context, team *Team) error
}

type TeamFilter struct {
	ID *int
}
