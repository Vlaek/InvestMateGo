package entity

type Marker interface {
	isRepository()
}

func (Asset) isRepository()    {}
func (Bond) isRepository()     {}
func (Share) isRepository()    {}
func (Etf) isRepository()      {}
func (Currency) isRepository() {}
