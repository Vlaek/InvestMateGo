package models

type RepositoryMarker interface {
	isRepository()
}

func (Bond) isRepository()  {}
func (Share) isRepository() {}
