package entity

type Marker interface {
	isRepository()
}

func (Bond) isRepository()  {}
func (Share) isRepository() {}
