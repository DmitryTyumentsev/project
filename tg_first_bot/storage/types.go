package storage

type Storage interface {
	IsExists(p *Page) (bool, error)
	PickRandom() (*Page, error)
	Save(p *Page) error
	Remove(p *Page) error
}

type Page struct {
	ChatID   int
	TextPage string
}

type InternalBasePath struct {
	BasePath string
}

func NewInternalBasePath(basePath string) *InternalBasePath {
	return &InternalBasePath{
		BasePath: basePath,
	}
}
