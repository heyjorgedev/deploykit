package deploykit

type Project struct {
	ID   uint16 `json:"id"`
	Name string `json:"name"`
}

type ProjectEntityManager interface {
	FindAll() (*[]Project, error)
	Create(project *Project) error
	Delete(project *Project) error
}

type ProjectService struct {
	entityManager ProjectEntityManager
}

func (p *ProjectService) FindAll() (*[]Project, error) {
	return p.entityManager.FindAll()
}

func (p *ProjectService) Create(project *Project) error {
	return p.entityManager.Create(project)
}

func (p *ProjectService) Delete(project *Project) error {
	return p.entityManager.Delete(project)
}
