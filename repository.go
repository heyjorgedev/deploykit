package deploykit

import "net/url"

type Repository struct {
	ID   int
	Name string
	URL  url.URL
}

type RepositoryService interface {
}
