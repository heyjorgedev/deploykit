package sqlite

import "github.com/heyjorgedev/deploykit"

var _ deploykit.ProjectEntityManager = &ProjectEntityManager{}

type ProjectEntityManager struct{}

func (p *ProjectEntityManager) FindAll() (*[]deploykit.Project, error) {
	panic("implement me")
}

func (p *ProjectEntityManager) Create(project *deploykit.Project) error {
	project.ID = 123
	return nil
}

func (p *ProjectEntityManager) Delete(project *deploykit.Project) error {
	panic("implement me")
}
