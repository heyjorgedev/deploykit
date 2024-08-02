package sqlite

import (
	"context"
	"strings"

	"github.com/heyjorgedev/deploykit"
)

type TeamService struct {
	db *DB
}

func NewTeamService(db *DB) *TeamService {
	return &TeamService{db: db}
}

func (s *TeamService) FindTeams(ctx context.Context, filter deploykit.TeamFilter) ([]*deploykit.Team, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()

	teams, n, err := findTeams(ctx, tx, filter)
	if err != nil {
		return teams, n, err
	}

	return teams, n, nil
}

func (s *TeamService) FindTeamByID(ctx context.Context, id int) (*deploykit.Team, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	team, err := findTeamByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (s *TeamService) CreateTeam(ctx context.Context, team *deploykit.Team) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createTeam(ctx, tx, team); err != nil {
		return err
	}

	return tx.Commit()
}

func findTeams(ctx context.Context, tx *Tx, filter deploykit.TeamFilter) (_ []*deploykit.Team, n int, err error) {
	where, args := []string{"1 = 1"}, []interface{}{}

	if filter.ID != nil {
		where = append(where, "id = ?")
		args = append(args, *filter.ID)
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    id,
		    name,
		    COUNT(*) OVER()
		FROM teams
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC`,
		args...,
	)
	if err != nil {
		return nil, n, formatError(err)
	}
	defer rows.Close()

	teams := make([]*deploykit.Team, 0)
	for rows.Next() {
		var team deploykit.Team
		if err := rows.Scan(
			&team.ID,
			&team.Name,
			&n,
		); err != nil {
			return nil, n, formatError(err)
		}

		teams = append(teams, &team)
	}

	return teams, 0, nil
}

func findTeamByID(ctx context.Context, tx *Tx, id int) (*deploykit.Team, error) {
	teams, _, err := findTeams(ctx, tx, deploykit.TeamFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(teams) == 0 {
		return nil, deploykit.Errorf(deploykit.ENOTFOUND, "Team not found.")
	}
	return teams[0], nil
}

func createTeam(ctx context.Context, tx *Tx, team *deploykit.Team) error {
	result, err := tx.ExecContext(ctx, "INSERT INTO teams (name) values (?);", team.Name)
	if err != nil {
		return formatError(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return formatError(err)
	}
	team.ID = int(id)
	return nil
}
