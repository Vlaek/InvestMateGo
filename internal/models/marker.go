package models

type DomainMarker interface {
	isDomain()
}

func (Bond) isDomain()     {}
func (Share) isDomain()    {}
func (Etf) isDomain()      {}
func (Currency) isDomain() {}
